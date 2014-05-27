package main

import (
    "log"
    "flag"
    "github.com/robmcl4/mycroft/srv"
    "github.com/robmcl4/mycroft/dispatch"
)

func main() {
    no_tls := flag.Bool("no-tls", false, "Whether to use TLS, default false")
    crt_path := flag.String("crt", "cert.crt", "Path to the TLS certificate, default `cert.crt`")
    key_path := flag.String("key", "key.key", "Path to the TLS key, default `key.key`")
    sname := flag.String("srv-name", "mycroft", "This server's name for SNI")

    flag.Parse()

    log.Println("Starting Dispatcher ...")
    go dispatch.Dispatch()

    log.Println("Starting Server ...")
    if *no_tls {
      log.Println("WARNING: not using TLS")
      err := srv.StartListen(1847, false, "", "", "")
      if err != nil {
        log.Println(err)
      }
    } else {
      err := srv.StartListen(1847, true, *crt_path, *key_path, *sname)
      if err != nil {
        log.Println(err)
      }
    }
}
