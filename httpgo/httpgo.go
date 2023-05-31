package main

import (
    "encoding/json"
    "fmt"
    "github.com/tarantool/go-tarantool"
    // "github.com/tarantool/go-tarantool/crud"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    // "reflect"
    "regexp"
    // "strings"
)

// func replace(w http.ResponseWriter, req *http.Request, space string, conn *tarantool.Connection) {
//     defer req.Body.Close()
//     body, err := ioutil.ReadAll(req.Body)

//     if err != nil {
//         http.Error(w, fmt.Sprintln("Failed to read body:", err), http.StatusInternalServerError)
//         return
//     }

//     data := []interface{}{}
//     err = json.Unmarshal(body, &data)

//     if err != nil {
//         http.Error(w, fmt.Sprintln("Failed to parse body:", err), http.StatusInternalServerError)
//         return
//     }

//     treq := crud.MakeReplaceRequest(space).Tuple(data)

//     type Tuple struct {
//         _msgpack struct{} `msgpack:",asArray"` //nolint: structcheck,unused
//         ID       uint64
//         BucketID uint64
//         Name     string
//     }
//     ret := crud.MakeResult(reflect.TypeOf(Tuple{}))

//     if err := conn.Do(treq).GetTyped(&ret); err != nil {
//         http.Error(w, fmt.Sprintln("Failed to execute request:", err), http.StatusInternalServerError)
//         return
//     }

//     fmt.Fprintln(w, ret)
// }

// var replaceRoute = regexp.MustCompile(`/data/\w+/replace`)

func main() {
    ConfigJSON := os.Getenv("TT_MICROSERVICE_CFG")

    type Config struct {
        Listen   string            `json:"listen"`
        Routes   map[string]string `json:"routes"`
        TTListen string            `json:"tt_listen"`
        TTUser   string            `json:"tt_user"`
        TTPass   string            `json:"tt_password"`
    }

    var cfg Config

    if err := json.Unmarshal([]byte(ConfigJSON), &cfg); err != nil {
        log.Fatalln("Failed to unmarshal config:", err)
    }

    opts := tarantool.Opts{User: cfg.TTUser, Pass: cfg.TTPass}
    conn, err := tarantool.Connect(cfg.TTListen, opts)
    if err != nil {
        log.Fatalln("Connection refused:", err)
    }

    httpRoutes := map[*regexp.Regexp]string{}

    for route, handler := range cfg.Routes {
        r, err := regexp.Compile(route)
        if err != nil {
            log.Fatalln("Failed to compile regex:", err)
        }
        httpRoutes[r] = handler
    }

    http.HandleFunc(
        "/",
        func(w http.ResponseWriter, req *http.Request) {
            for route, handler := range httpRoutes {
                if route.Match([]byte(req.URL.Path)) {
                    defer req.Body.Close()
                    body, err := ioutil.ReadAll(req.Body)

                    if err != nil {
                        http.Error(w, fmt.Sprintln("Failed to read body:", err), http.StatusInternalServerError)
                        return
                    }

                    call := tarantool.NewCallRequest(handler).Args([]interface{}{body})
                    resp, err := conn.Do(call).Get()

                    if err != nil {
                        http.Error(w, fmt.Sprintln("Failed to execute request:", err), http.StatusInternalServerError)
                        return
                    }

                    if len(resp.Data) == 0 {
                        fmt.Fprintln(w, "empty response")
                    } else {
                        fmt.Fprintln(w, resp.Data[0])
                    }

                    return
                }
            }

            http.NotFound(w, req)
        },
    )

    if err := http.ListenAndServe(cfg.Listen, nil); err != nil {
        log.Fatalln("Failed to start a server:", err)
    }
}
