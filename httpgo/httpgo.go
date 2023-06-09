package main

import (
    "encoding/json"
    "github.com/tarantool/go-tarantool"
    "log"
    "net/http"
    "os"
    "plugin"
    // "regexp"
)

func main() {
    // Prepare Tarantool connection
    ServerListen := os.Getenv("SERVER_LISTEN")
    ServerUser := os.Getenv("SERVER_USER")
    ServerPass := os.Getenv("SERVER_PASS")

    opts := tarantool.Opts{User: ServerUser, Pass: ServerPass}
    conn, err := tarantool.Connect(ServerListen, opts)
    if err != nil {
        log.Fatalln("Connection refused:", err)
    }

    // Parse server + plugins configuration
    ConfigJSON := os.Getenv("TT_MICROSERVICE_CFG")
    type Config struct {
        Listen   string `json:"listen"`
        Pipeline []struct {
            Plugin struct {
                Lib    string `json:"lib"`
                Symbol string `json:"symbol"`
            } `json:"plugin"`
            Cfg map[string]interface{} `json:"cfg"`
        } `json:"pipeline"`
    }

    var cfg Config

    if err := json.Unmarshal([]byte(ConfigJSON), &cfg); err != nil {
        log.Fatalln("Failed to unmarshal config:", err)
    }

    // Prepare route processor. The last processor is default http.NotFound
    handler := http.NotFound
    for i := len(cfg.Pipeline) - 1; i >= 0; i-- {
        component := cfg.Pipeline[i]
        p, err := plugin.Open(component.Plugin.Lib)
        if err != nil {
            log.Fatalf("Failed to open plugin %s: %w", component.Plugin.Lib, err)
        }

        symbol, err := p.Lookup(component.Plugin.Symbol)
        if err != nil {
            log.Fatalf("Failed to open lookup plugin %s symbol %s : %w",
                component.Plugin.Lib, component.Plugin.Symbol, err)
        }

        if prevHandlerRaw, ok := symbol.(func(http.ResponseWriter,
            *http.Request,
            *tarantool.Connection,
            map[string]interface{},
            func(http.ResponseWriter, *http.Request),
        )); !ok {
            log.Fatalf("Failed to cast plugin %s symbol %s to expected type", component.Plugin.Lib, component.Plugin.Symbol)
        } else {
            currentHandler := handler
            handler = func(w http.ResponseWriter, req *http.Request) {
                prevHandlerRaw(w, req, conn, component.Cfg, currentHandler)
            }
        }
    }

    http.HandleFunc("/", handler)

    if err := http.ListenAndServe(cfg.Listen, nil); err != nil {
        log.Fatalln("Failed to start a server:", err)
    }
}
