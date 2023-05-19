package main

import (
    "github.com/tarantool/go-tarantool"
    "log"
    "os"
    // "github.com/tarantool/go-tarantool/crud"
)

// type person struct {
//     name string
//     age  int
// }

// func server_start()

func main() {
    ListenAddr := os.Getenv("LISTEN_ADDR")
    ServerUser := os.Getenv("SERVER_USER")
    ServerPass := os.Getenv("SERVER_PASS")

    opts := tarantool.Opts{User: ServerUser, Pass: ServerPass}
    conn, err := tarantool.Connect(ListenAddr, opts)
    if err != nil {
        log.Fatalln("Connection refused:", err)
    }

    resp, err := conn.Call("get_server_config", []interface{}{})
    if err != nil {
        log.Fatalln("Failed:", err)
    }

    log.Println(resp)
}
