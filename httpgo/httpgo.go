package main

import (
    "encoding/json"
    "fmt"
    "github.com/tarantool/go-tarantool"
    // "github.com/tarantool/go-tarantool/crud"
    // "io/ioutil"
    "log"
    "net/http"
    "os"
    // "reflect"
    // "regexp"
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
    }

    var cfg Config

    if err := json.Unmarshal([]byte(ConfigJSON), &cfg); err != nil {
        log.Fatalln("Failed to unmarshal config:", err)
    }

    http.HandleFunc(
        "/hello",
        func(w http.ResponseWriter, req *http.Request) {
            eval := tarantool.NewEvalRequest("return 'Hello world!'")
            resp, err := conn.Do(eval).Get()
            if err != nil {
                http.Error(w, fmt.Sprintln("Failed to execute request:", err), http.StatusInternalServerError)
                return
            }

            if len(resp.Data) == 0 {
                http.Error(w, "Got unexpected zero length response", http.StatusInternalServerError)
                return
            }

            // switch {
            // case replaceRoute.MatchString(req.URL.Path):
            //     p := strings.Split(req.URL.Path, "/")
            //     replace(w, req, p[2], conn)
            // default:
            //     http.NotFound(w, req)
            // }
            fmt.Fprintln(w, resp.Data[0])
        },
    )

    if err := http.ListenAndServe(cfg.Listen, nil); err != nil {
        log.Fatalln("Failed to start a server:", err)
    }
}
