package app


import (
    "net"
    "encoding/json"
    "fmt"
    "log"
    "github.com/coreos/go-semver/semver"
    "sync"
)


const (
    STATUS_CONNECTED = iota
    STATUS_UP = iota
    STATUS_DOWN = iota
    STATUS_IN_USE = iota
)


type Capability struct {
    Version semver.Version
    Name string
}


type Manifest struct {
    Name string
    DisplayName string
    InstanceId string
    ApiVersion int
    Description string
    Version semver.Version
    Capabilities []*Capability
    Dependencies []*Capability
}


type App struct {
    Connection net.Conn
    Manifest *Manifest
    Status int
    Priority int
    RWMutex sync.RWMutex
}


func NewApp() (*App) {
    ret := new(App)
    ret.Status = STATUS_CONNECTED
    return ret
}


func (a *App) Send(verb string, body map[string]interface{}) (error) {
    id := ""
    if a.Manifest != nil {
        id = a.Manifest.InstanceId
    }
    log.Printf("Sending message %s to instance %s\n", verb, id)
    verb += " "
    bodyBytes, err := json.Marshal(body)
    if err != nil {
        return err
    }
    toSend := append([]byte(verb), bodyBytes...)
    length := fmt.Sprintf("%d\n", len(toSend))
    toSend = append([]byte(length), toSend...)
    _, err = a.Connection.Write(toSend)
    return err
}


func (a *App) StatusString() (string) {
    switch a.Status {
    case STATUS_CONNECTED:
        return "connected"
    case STATUS_UP:
        return "up"
    case STATUS_IN_USE:
        return fmt.Sprintf("in_use %d", a.Priority)
    case STATUS_DOWN:
        return "down"
    default:
        return ""
    }
}
