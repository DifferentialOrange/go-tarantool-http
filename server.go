package main

import (
    "github.com/tarantool/go-tarantool"
    // "github.com/tarantool/go-tarantool/crud"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
)

// func init_server() {
//     resp, err := conn.Call("get_server_config", []interface{}{})
//     if err != nil {
//         log.Fatalln("Failed:", err)
//     }

//     // go server_start(cfg)
//     go server_start()
// }

// //  func server_start(cfg []interface{}) {
// //      log.Println(cfg)
// //  }
// func server_start() {

// }

func main() {
    ListenAddr := os.Getenv("LISTEN_ADDR")
    ServerUser := os.Getenv("SERVER_USER")
    ServerPass := os.Getenv("SERVER_PASS")

    opts := tarantool.Opts{User: ServerUser, Pass: ServerPass}
    conn, err := tarantool.Connect(ListenAddr, opts)
    if err != nil {
        log.Fatalln("Connection refused:", err)
    }

    mux := http.NewServeMux()

    mux.HandleFunc(
        "/data/insert",
        func(w http.ResponseWriter, req *http.Request) {
            resp, err := conn.Eval("return true", []interface{}{})
            if err != nil {
                log.Fatalln("Eval failed:", err)
            }
            fmt.Fprintln(w, "Welcome to the home page! We've got", resp)
        },
    )

    s := &http.Server{
        Addr:           "localhost:8080",
        Handler:        mux,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }
    log.Fatal(s.ListenAndServe())
}
