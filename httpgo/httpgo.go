package main

import (
    "encoding/json"
    // "fmt"
    "github.com/tarantool/go-tarantool"
    // "github.com/tarantool/go-tarantool/crud"
    // "io/ioutil"
    "log"
    "net/http"
    "os"
    // "reflect"
    "regexp"
    // "strings"
    // "path/filepath"
    "plugin"
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
    ServerListen := os.Getenv("SERVER_LISTEN")
    ServerUser := os.Getenv("SERVER_USER")
    ServerPass := os.Getenv("SERVER_PASS")

    ConfigJSON := os.Getenv("TT_MICROSERVICE_CFG")

    opts := tarantool.Opts{User: ServerUser, Pass: ServerPass}
    conn, err := tarantool.Connect(ServerListen, opts)
    if err != nil {
        log.Fatalln("Connection refused:", err)
    }

    type Config struct {
        Listen string `json:"listen"`
        Routes []struct {
            Handler struct {
                Type string `json:"type"`
                Lib  string `json:"lib"`
                Func string `json:"func"`
            } `json:"handler"`
            Path string `json:"path"`
        } `json:"routes"`
    }

    var cfg Config

    if err := json.Unmarshal([]byte(ConfigJSON), &cfg); err != nil {
        log.Fatalln("Failed to unmarshal config:", err)
    }

    httpRoutes := map[*regexp.Regexp]func(http.ResponseWriter, *http.Request, *tarantool.Connection){}

    for _, route := range cfg.Routes {
        r, err := regexp.Compile(route.Path)
        if err != nil {
            log.Fatalln("Failed to compile regex:", err)
        }

        if route.Handler.Type == "plugin" {
            p, err := plugin.Open(route.Handler.Lib)
            if err != nil {
                log.Fatalf("Failed to open plugin %s: %w", route.Handler.Lib, err)
            }

            v, err := p.Lookup(route.Handler.Func)
            if err != nil {
                log.Fatalf("Failed to open lookup plugin %s function %s : %w", route.Handler.Lib, route.Handler.Func, err)
            }

            if handler, ok := v.(func(http.ResponseWriter, *http.Request, *tarantool.Connection)); !ok {
                log.Fatalf("Failed to cast plugin %s function %s to expected type", route.Handler.Lib, route.Handler.Func)
            } else {
                httpRoutes[r] = handler
            }
        } else {
            log.Fatalf("Failed to start server: unknown routre source type '%s'", route.Handler.Type)
        }
    }

    http.HandleFunc(
        "/",
        func(w http.ResponseWriter, req *http.Request) {
            for route, handler := range httpRoutes {
                if route.Match([]byte(req.URL.Path)) {
                    handler(w, req, conn)
                    return
                    // defer req.Body.Close()
                    // body, err := ioutil.ReadAll(req.Body)

                    // if err != nil {
                    //     http.Error(w, fmt.Sprintln("Failed to read body:", err), http.StatusInternalServerError)
                    //     return
                    // }

                    // call := tarantool.NewCallRequest(handler).Args([]interface{}{body})
                    // resp, err := conn.Do(call).Get()

                    // if err != nil {
                    //     http.Error(w, fmt.Sprintln("Failed to execute request:", err), http.StatusInternalServerError)
                    //     return
                    // }

                    // if len(resp.Data) == 0 {
                    //     fmt.Fprintln(w, "empty response")
                    // } else {
                    //     fmt.Fprintln(w, resp.Data[0])
                    // }

                    // return
                }
            }

            http.NotFound(w, req)
        },
    )

    if err := http.ListenAndServe(cfg.Listen, nil); err != nil {
        log.Fatalln("Failed to start a server:", err)
    }
}
