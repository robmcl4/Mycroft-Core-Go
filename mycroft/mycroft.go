// Package main contains the main entry point for the Mycroft server.
package main

import (
    "flag"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/srv"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/logging"
)


// main parses command-line arguments and spawns a new server
func main() {
    logging.SetLevel(logging.DEBUG)

    no_tls := flag.Bool("no-tls", false, "Whether to use TLS, default false")
    crt_path := flag.String("crt", "cert.crt", "Path to the TLS certificate, default `cert.crt`")
    key_path := flag.String("key", "key.key", "Path to the TLS key, default `key.key`")
    sname := flag.String("srv-name", "mycroft", "This server's name for SNI")

    flag.Parse()

    logging.Info("Starting Server ...")

    if *no_tls {
        logging.Warning("not using TLS")
        err := srv.StartListen(1847, false, "", "", "")
        if err != nil {
            logging.Fatal("Could not start server: ", err.Error())
        }
    } else {
        err := srv.StartListen(1847, true, *crt_path, *key_path, *sname)
        if err != nil {
            logging.Fatal("Could not start server: ", err.Error())
        }
    }
}
