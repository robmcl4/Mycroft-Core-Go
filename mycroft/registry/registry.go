package registry

import (
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/registry/msg_archive"
)

var capabilitySuppliers map[app.Capability] []*app.App = make(map[app.Capability][]*app.App)
var capabilityDependents map[app.Capability] []*app.App = make(map[app.Capability][]*app.App)
var instances map[string] *app.App = make(map[string]*app.App)


// inserts into a given map so map[capability] points to the given app
func addCapabilityToMap(m map[app.Capability] []*app.App, a *app.App, cpb *app.Capability) {
    m[*cpb] = append(m[*cpb], a)
}

// removes an app's capability from the given map and returns true if removed
func removeCapabilityFromMap(m map[app.Capability] []*app.App, a *app.App, cpb *app.Capability) {
    if apps, ok := m[*cpb]; ok {
        for i, curr := range apps {
            if a == curr {
                // check out https://code.google.com/p/go-wiki/wiki/SliceTricks
                // for an explanation
                copy(apps[i:], apps[i+1:])
                apps[len(apps)-1] = nil
                m[*cpb] = apps[:len(apps)-1]
                break
            }
        }
    }
}


// adds all of an app's capabilities to the registry
func addCapabilities(a *app.App) {
    for _, val := range a.Manifest.Capabilities {
        addCapabilityToMap(capabilitySuppliers, a, val)
    }
}


// removes all of an app's capabilities from the registry
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


// adds an app's dependencies to the registry
func removeDependencies(a *app.App) {
    for _, val := range a.Manifest.Dependencies {
        removeCapabilityFromMap(capabilityDependents, a, val)
    }
}


// adds an app's instance id to the registry
func addInstanceId(a *app.App) {
    instances[a.Manifest.InstanceId] = a
}


// removes an app's instance id from the registry
func removeInstanceId(a *app.App) {
    delete(instances, a.Manifest.InstanceId)
}


// Registers an app with the registry. This is currently NOT thread-safe.
func Register(a *app.App) {
    addCapabilities(a)
    addDependencies(a)
    addInstanceId(a)
}


// Removes an app from the registry. This is currently NOT thread-safe.
func Remove(a *app.App) {
    removeCapabilities(a)
    removeDependencies(a)
    removeInstanceId(a)
    msg_archive.RemoveAppsMessages(a)
}


// Gets a list of all apps that provide this capability
func GetProviders(cpb *app.Capability) ([]*app.App) {
    if apps, ok := capabilitySuppliers[*cpb]; ok {
        ret := make([]*app.App, len(apps))
        copy(ret, apps)
        return ret
    }
    return make([]*app.App, 0)
}


// Gets a list of all apps that depend on this capability
func GetDependents(cpb *app.Capability) ([]*app.App) {
    if apps, ok := capabilityDependents[*cpb]; ok {
        ret := make([]*app.App, len(apps))
        copy(ret, apps)
        return ret
    }
    return make([]*app.App, 0)
}


// Gets the app known by the given instance ID
func GetInstance(instId string) (a *app.App, ok bool) {
    a, ok = instances[instId]
    return
}
