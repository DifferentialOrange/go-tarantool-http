package main

import (
    // "encoding/json"
    "github.com/tarantool/go-tarantool"
    "github.com/urfave/negroni"
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
    "log"
    "strconv"
    "time"
)

func Measure(conn *tarantool.Connection, collectorName string,
    elapsed time.Duration, path string, method string, status int,
) {
    eval := tarantool.NewEvalRequest(`
local metrics = require('metrics')

local collector_name, time, path, method, status = ...

metrics.summary(collector_name):observe(time, {path = path, method = method, status = status})
`).Args([]interface{}{collectorName, elapsed.Seconds(), path, method, strconv.Itoa(status)})
    _, err := conn.Do(eval).Get()
    if err != nil {
        log.Println("Failed to send metrics:", err)
    }
}

func MainHandler(w http.ResponseWriter,
    req *http.Request,
    conn *tarantool.Connection,
    cfg map[string]interface{},
    nextHandler func(http.ResponseWriter, *http.Request),
) {
    collectorNameRaw, ok := cfg["collector_name"]
    if !ok {
        http.Error(w, "Got unexpected config: no collector_name", http.StatusInternalServerError)
        return
    }

    collectorName, ok := collectorNameRaw.(string)
    if !ok {
        http.Error(w, "Failed to cast collector_name to string", http.StatusInternalServerError)
        return
    }

    start := time.Now()

    lrw := negroni.NewResponseWriter(w)

    nextHandler(lrw, req)

    elapsed := time.Since(start)

    go Measure(conn, collectorName, elapsed, req.URL.Path, req.Method, lrw.Status())
}
