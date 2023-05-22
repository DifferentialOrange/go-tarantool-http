package main

import (
    "encoding/json"
    "fmt"
    "github.com/tarantool/go-tarantool"
    "github.com/tarantool/go-tarantool/crud"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "reflect"
    "regexp"
    "strings"
)

func replace(w http.ResponseWriter, req *http.Request, space string, conn *tarantool.Connection) {
    defer req.Body.Close()
    body, err := ioutil.ReadAll(req.Body)

    if err != nil {
        http.Error(w, fmt.Sprintln("Failed to read body:", err), http.StatusInternalServerError)
        return
    }

    data := []interface{}{}
    err = json.Unmarshal(body, &data)

    if err != nil {
        http.Error(w, fmt.Sprintln("Failed to parse body:", err), http.StatusInternalServerError)
        return
    }

    treq := crud.MakeReplaceRequest(space).Tuple(data)

    type Tuple struct {
        _msgpack struct{} `msgpack:",asArray"` //nolint: structcheck,unused
        ID       uint64
        BucketID uint64
        Name     string
    }
    ret := crud.MakeResult(reflect.TypeOf(Tuple{}))

    if err := conn.Do(treq).GetTyped(&ret); err != nil {
        http.Error(w, fmt.Sprintln("Failed to execute request:", err), http.StatusInternalServerError)
        return
    }

    fmt.Fprintln(w, ret)
}

var replaceRoute = regexp.MustCompile(`/data/\w+/replace`)

func main() {
    ListenAddr := os.Getenv("LISTEN_ADDR")
    ServerUser := os.Getenv("SERVER_USER")
    ServerPass := os.Getenv("SERVER_PASS")

    opts := tarantool.Opts{User: ServerUser, Pass: ServerPass}
    conn, err := tarantool.Connect(ListenAddr, opts)
    if err != nil {
        log.Fatalln("Connection refused:", err)
    }

    http.HandleFunc(
        "/",
        func(w http.ResponseWriter, req *http.Request) {
            switch {
            case replaceRoute.MatchString(req.URL.Path):
                p := strings.Split(req.URL.Path, "/")
                replace(w, req, p[2], conn)
            default:
                http.NotFound(w, req)
            }
        },
    )

    http.ListenAndServe(":8080", nil)
}
