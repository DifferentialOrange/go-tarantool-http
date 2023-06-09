package main

import (
    "github.com/tarantool/go-tarantool"
    "net/http"
    "strings"
)

func MainHandler(w http.ResponseWriter,
    req *http.Request,
    conn *tarantool.Connection,
    cfg map[string]interface{},
    nextHandler func(http.ResponseWriter, *http.Request),
) {
    expectedTokenRaw, ok := cfg["token"]
    if !ok {
        http.Error(w, "Got unexpected config: no token", http.StatusInternalServerError)
        return
    }

    expectedToken, ok := expectedTokenRaw.(string)
    if !ok {
        http.Error(w, "Failed to cast token to string", http.StatusInternalServerError)
        return
    }

    reqToken := req.Header.Get("Authorization")
    if reqToken == "" {
        http.Error(w, "No Authorization token", http.StatusUnauthorized)
        return
    }
    splitToken := strings.Split(reqToken, "Basic ")
    if len(splitToken) < 2 {
        http.Error(w, "No Basic Authorization token", http.StatusUnauthorized)
        return
    }
    token := splitToken[1]

    if token != expectedToken {
        http.Error(w, "Wrong token", http.StatusUnauthorized)
        return
    }

    // Pass to next handler only if request wan't processed by any handler here.
    nextHandler(w, req)
}
