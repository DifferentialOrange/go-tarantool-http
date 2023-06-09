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
    "regexp"
)

var Handlers map[string]func(http.ResponseWriter, *http.Request, *tarantool.Connection) = map[string]func(http.ResponseWriter, *http.Request, *tarantool.Connection){
    "Hello": func(w http.ResponseWriter, req *http.Request, conn *tarantool.Connection) {
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
    },
}

func MainHandler(w http.ResponseWriter,
    req *http.Request,
    conn *tarantool.Connection,
    cfg map[string]interface{},
    nextHandler func(http.ResponseWriter, *http.Request)) {
    routesRaw, ok := cfg["routes"]
    if !ok {
        http.Error(w, "Got unexpected config: no routes", http.StatusInternalServerError)
        return
    }

    routes, ok := routesRaw.([]interface{})
    if !ok {
        http.Error(w, "Failed to cast routes to []interface{}", http.StatusInternalServerError)
        return
    }

    for _, routeRaw := range routes {
        route, ok := routeRaw.(map[string]interface{})
        if !ok {
            http.Error(w, "Failed to cast route to map[string]interface", http.StatusInternalServerError)
            return
        }

        methodRaw, ok := route["method"]
        if !ok {
            http.Error(w, "Missing method in route", http.StatusInternalServerError)
            return
        }
        method, ok := methodRaw.(string)
        if !ok {
            http.Error(w, "Failed to parse method to string", http.StatusInternalServerError)
            return
        }

        pathRaw, ok := route["path"]
        if !ok {
            http.Error(w, "Missing path in route", http.StatusInternalServerError)
            return
        }
        path, ok := pathRaw.(string)
        if !ok {
            http.Error(w, "Failed to parse path to string", http.StatusInternalServerError)
            return
        }

        handlerRaw, ok := route["handler"]
        if !ok {
            http.Error(w, "Missing handler in route", http.StatusInternalServerError)
            return
        }
        handler, ok := handlerRaw.(string)
        if !ok {
            http.Error(w, "Failed to parse handler to string", http.StatusInternalServerError)
            return
        }

        r, err := regexp.Compile(path)
        if !ok {
            http.Error(w, fmt.Sprintln("Failed to compile a path:", err), http.StatusInternalServerError)
            return
        }

        if r.Match([]byte(req.URL.Path)) && (req.Method == method) {
            handlerFunc, ok := Handlers[handler]
            if !ok {
                http.Error(w, fmt.Sprintln("Handler unknown:", handler), http.StatusInternalServerError)
                return
            }

            handlerFunc(w, req, conn)
            return
        }
    }

    // Pass to next handler only if request wan't processed by any handler here.
    nextHandler(w, req)
}
