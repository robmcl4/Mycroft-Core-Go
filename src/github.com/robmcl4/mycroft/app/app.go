package app


import (
    "net"
    "github.com/coreos/go-semver/semver"
)


const (
    STATUS_CONNECTED = iota
    STATUS_UP = iota
    STATUS_DOWN = iota
    STATUS_IN_USE = iota
)


type Capability struct {
    Version *semver.Version
    Name string
}


type Manifest struct {
    Name string
    DisplayName string
    InstanceId string
    ApiVersion int
    Description string
    Version string
    Capabilities []*Capability
    Dependencies []*Capability
}


type App struct {
    Connection *net.Conn
    Manifest *Manifest
    Status int
}


func NewApp() (*App) {
    ret := new(App)
    ret.Status = STATUS_CONNECTED
    return ret
}
