package registry

import (
    "github.com/robmcl4/mycroft/app"
    "github.com/coreos/go-semver/semver"
)


var capabilitySuppliers map[string] map[semver.Version] []*app.App = make(map[string]map[semver.Version][]*app.App)
var capabilityDependents map[string] map[semver.Version] []*app.App = make(map[string]map[semver.Version][]*app.App)
var instances map[string] *app.App = make(map[string]*app.App)


// insert into a given map so map[name][version] points to the given app
func addCapabilityToMap(m map[string] map[semver.Version] []*app.App, a *app.App, cpb *app.Capability) {
    name := cpb.Name
    ver := *cpb.Version
    // add the version map to capabilitySuppliers if needed
    if _, ok := m[name]; !ok {
        m[name] = make(map[semver.Version][]*app.App)
    }
    // add the array to capabilitySuppliers if needed
    if _, ok := m[name][ver]; !ok {
        m[name][ver] = make([]*app.App, 0)
    }
    m[name][ver] = append(m[name][ver], []*app.App{a}...)
}


// add an app's capabilities to the registry
func addCapabilities(a *app.App) {
    for _, val := range a.Manifest.Capabilities {
        addCapabilityToMap(capabilitySuppliers, a, val)
    }
}


// add an app's dependencies to the registry
func addDependencies(a *app.App) {
    for _, val := range a.Manifest.Dependencies {
        addCapabilityToMap(capabilityDependents, a, val)
    }
}


// add an app's instance id to the registry
func addInstanceId(a *app.App) {
    instances[a.Manifest.InstanceId] = a
}


func Register(a *app.App) {
    addCapabilities(a)
    addDependencies(a)
    addInstanceId(a)
}


// Get apps that provide this capability
func GetProviders(cpb *app.Capability) ([]*app.App) {
    ret := make([]*app.App, 0)
    if _, ok := capabilitySuppliers[cpb.Name]; ok {
        if apps, ok := capabilitySuppliers[cpb.Name][*cpb.Version]; ok {
            ret := make([]*app.App, len(apps))
            copy(ret, apps)
            return ret
        }
    }
    return ret
}


func GetDependents(cpb *app.Capability) ([]*app.App) {
    ret := make([]*app.App, 0)
    if _, ok := capabilityDependents[cpb.Name]; ok {
        if apps, ok := capabilityDependents[cpb.Name][*cpb.Version]; ok {
            ret := make([]*app.App, len(apps))
            copy(ret, apps)
            return ret
        }
    }
    return ret
}


func GetInstance(instId string) (a *app.App, ok bool) {
    a, ok = instances[instId]
    return
}
