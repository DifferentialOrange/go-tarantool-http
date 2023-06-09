package main

import (
    // "encoding/json"
    "fmt"
    "github.com/tarantool/go-tarantool"
    // "github.com/tarantool/go-tarantool/crud"
    // "io/ioutil"
    // "log"
    "net/http"
    // "os"
    // "reflect"
    // "regexp"
    // "strings"
    // "path/filepath"
    // "plugin"
)

func init() {
    // Name := "httpgo-crud"
}

func main() {
    // Name := "httpgo-crud"
}

func Hello(w http.ResponseWriter, req *http.Request, conn *tarantool.Connection) {
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

    fmt.Fprintln(w, resp.Data[0])
}
