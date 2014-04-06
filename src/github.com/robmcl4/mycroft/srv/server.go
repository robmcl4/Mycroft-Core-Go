package srv

import (
    "net"
    "log"
    "fmt"
    "strings"
    "strconv"
    "github.com/robmcl4/mycroft/app"
    "github.com/robmcl4/mycroft/cmd"
    "github.com/robmcl4/mycroft/dispatch"
    "github.com/robmcl4/mycroft/registry"
)


// Start listening for applications on the given port.
// When new applications connect it will launch them in a new gorutine.
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
    defer l.Close() // at the end of this method close the connection
    log.Println("Starting listen loop")
    for {
        a, err := acceptApp(l)
        if err != nil {
            log.Fatal(err)
            return err
        } else {
            log.Println("Got connection")
            go ListenForCommands(a)
        }
    }
    return nil
}


// Listen for and accept a new application connection
func acceptApp(lnr *net.TCPListener) (*app.App, error) {
    conn, err := lnr.Accept()
    if err != nil {
        return nil, err
    }
    ret := app.NewApp()
    ret.Connection = conn
    return ret, nil
}


// Start listening for commands through this app's connection.
// NOTE: this should likely only be called as a goroutine.
func ListenForCommands(a *app.App) {
    defer closeApp(a)
    smallBuff := make([]byte, 200)
    smallBuffI := 0
    for smallBuffI < len(smallBuff) {
        innerBuff := make([]byte, 1)
        _, err := a.Connection.Read(innerBuff)
        if err != nil {
            log.Println("ERROR:", err)
            return
        }
        smallBuff[smallBuffI] = innerBuff[0]
        smallBuffI += 1
        str := string(smallBuff[:smallBuffI])
        if len(str) > 0 && strings.HasSuffix(str, "\n") {
            // whoa we found a message length! read it
            var msgLen int64
            msgLen, err = strconv.ParseInt(str[:len(str)-1], 10, 64)
            if err != nil {
                log.Printf("ERROR: could not parse '%s': %s\n", str, err.Error())
                return
            }
            msgBuff := make([]byte, msgLen)
            n, err := a.Connection.Read(msgBuff)
            if err != nil {
                log.Println("ERROR:", err)
                return
            }
            cmd := cmd.ParseCommand(a, msgBuff[:n])
            if err != nil {
                log.Println("ERROR:", err)
                return
            }
            dispatch.Enqueue(cmd)
            smallBuff = make([]byte, 200)
            smallBuffI = 0
        }
    }
    log.Printf("Closing connection to app, garbage was read")
}


// perform all operations required to close this app
// this really should be somewhere else in the code, but i can't figure out where
// since most places would lead to circular references
func closeApp(a *app.App) {
    a.Connection.Close()
    registry.Remove(a)
    sc, _ := cmd.NewStatusChange(a, app.STATUS_DOWN, nil)
    dispatch.Enqueue(sc)
    if a.Manifest != nil {
        log.Printf("Closing application %s", a.Manifest.InstanceId)
    } else {
        log.Printf("Closing application")
    }
}
