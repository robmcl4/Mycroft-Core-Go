package srv

import (
    "net"
    "log"
    "fmt"
    "github.com/robmcl4/mycroft/app"
)


// start listening for applications on the given port
func StartListen(port int) (error) {
    addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        log.Fatal(err)
        return err
    }
    l, err := net.ListenTCP("tcp", addr)
    if err != nil {
        log.Fatal(err)
        return err
    }
    defer l.Close()
    log.Println("Starting listen loop")
    for {
        _, err := acceptApp(l)
        if err != nil {
            log.Fatal(err)
        } else {
            log.Println("Got connection")
        }
    }
    return nil
}


// listen for and accept a new application connection
func acceptApp(lnr *net.TCPListener) (*app.App, error) {
    conn, err := lnr.Accept()
    if err != nil {
        return nil, err
    }
    ret := app.NewApp()
    ret.Connection = &conn
    return ret, nil
}