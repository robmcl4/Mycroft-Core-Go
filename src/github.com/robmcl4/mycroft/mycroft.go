package main

import (
    "log"
    "github.com/robmcl4/mycroft/srv"
    "github.com/robmcl4/mycroft/dispatch"
)

func main() {
    log.Println("Starting Server ...")
    go dispatch.Dispatch()
    srv.StartListen(1847)
}
