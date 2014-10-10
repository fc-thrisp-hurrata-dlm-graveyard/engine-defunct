+++
title = "Install & Quick Start"
+++
# Install & QuickStart

### Install

    go get github.com/thrisp/engine

### Quickstart

main.go

    package main

    import (
        "os"
        "os/signal"
        "github.com/thrisp/engine"
    )

    func Display(c *engine.Ctx) {
        c.RW.Header().Set("Content-Type", "text/html")
        c.RW.Write("HELLO WORLD!")
    }

    var quit = make(chan bool)

    func init() {
        c := make(chan os.Signal, 1)
        signal.Notify(c, os.Interrupt)
        go func() {
            for _ = range c {
                quit <- true
            }
        }()
    }

    func main() {
        e := engine.Basic()
        e.Manage("GET", "/hello/world/", Display)
        go e.Run(":8080")
        <-quit
    }

go run main.go
