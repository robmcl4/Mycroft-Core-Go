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


func removeCapabilityFromMap(m map[string] map[semver.Version] []*app.App, a *app.App, cpb *app.Capability) (bool) {
    name := cpb.Name
    ver := *cpb.Version
    // make sure the capability name exists
    if _, ok := m[name]; !ok {
        return false
    }
    // make sure the capability version exists
    if apps, ok := m[name][ver]; !ok {
        return false
    } else {
        // both exist, there must be a non-nil array in val
        // see if we can find this app in that array, then store an array
        // without the app
        for i, app := range apps {
            if app == a {
                // i don't know either. 
                // check out https://code.google.com/p/go-wiki/wiki/SliceTricks
                copy(apps[i:], apps[i+1:])
                apps[len(apps)-1] = nil
                apps = apps[:len(apps)-1]
                m[name][ver] = apps
                return true
            }
        }
    }
    return false
}


// add an app's capabilities to the registry
func addCapabilities(a *app.App) {
    for _, val := range a.Manifest.Capabilities {
        addCapabilityToMap(capabilitySuppliers, a, val)
    }
}


func removeCapabilities(a *app.App) {
    for _, val := range a.Manifest.Capabilities {
        removeCapabilityFromMap(capabilitySuppliers, a, val)
    }
}


// add an app's dependencies to the registry
func addDependencies(a *app.App) {
    for _, val := range a.Manifest.Dependencies {
        addCapabilityToMap(capabilityDependents, a, val)
    }
}


func removeDependencies(a *app.App) {
    for _, val := range a.Manifest.Dependencies {
        removeCapabilityFromMap(capabilityDependents, a, val)
    }
}


// add an app's instance id to the registry
func addInstanceId(a *app.App) {
    instances[a.Manifest.InstanceId] = a
}


func removeInstanceId(a *app.App) {
    delete(instances, a.Manifest.InstanceId)
}


func Register(a *app.App) {
    addCapabilities(a)
    addDependencies(a)
    addInstanceId(a)
}


func Remove(a *app.App) {
    removeCapabilities(a)
    removeDependencies(a)
    removeInstanceId(a)
    removeAppsMessages(a)
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
