package cmd

import (
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/registry"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/logging"
    "github.com/coreos/go-semver/semver"
    "github.com/nu7hatch/gouuid"
    "errors"
    "fmt"
)

// construct the app based on its manifest
func (c *commandStrategy) appManifest() (bool) {
    c.app.RWMutex.Lock()
    defer c.app.RWMutex.Unlock()

    logging.Debug("Parsing application's manifest")

    man, err := decodeManifest(c.body)
    if err != nil {
        logging.Error("App's manifest is invalid: %s", err)
        sendManifestFail(c.app, err.Error())
        return false
    }
    if _, exists := registry.GetInstance(man.InstanceId); exists {
        logging.Error("App's instance id is already taken: %s", man.InstanceId)
        sendManifestFail(c.app, fmt.Sprintf("instanceId %s is in use", man.InstanceId))
        return false
    }
    c.app.Status = app.STATUS_DOWN
    c.app.Manifest = man
    registry.Register(c.app)
    sendManifestOkAndDependencies(c.app)
    logging.Info("App '%s' now connected with manifest parsed", c.app.Manifest.InstanceId)
    return true
}


func sendManifestFail(a *app.App, reason string) {
    body := make(map[string]interface{})
    body["message"] = reason
    a.Send("APP_MANIFEST_FAIL", body)
}


func sendManifestOkAndDependencies(a *app.App) {
    // send manifest ok
    body := make(map[string]interface{})
    body["instanceId"] = a.Manifest.InstanceId
    a.Send("APP_MANIFEST_OK", body)
    // send dependencies
    body = make(map[string]interface{})
    for _, dep := range a.Manifest.Dependencies {
        inner := make(map[string] string)
        body[dep.Name] = inner
        for _, provider := range registry.GetProviders(dep) {
            inner[provider.Manifest.InstanceId] = provider.StatusString()
        }
    }
    a.Send("APP_DEPENDENCY", body)
}


func decodeManifest(m jsonData) (man *app.Manifest, err error) {
    man = new(app.Manifest)

    // start loading in values from the manifest
    if str, ok := getString(m, "name"); ok && len(str) != 0 {
        man.Name = str
    } else {
        err = errors.New("No name was found")
        return
    }

    if val, ok := getString(m, "displayName"); ok && len(val) != 0 {
        man.DisplayName = val
    } else {
        err = errors.New("No displayName was found")
        return
    }

    if val, ok := getString(m, "instanceId"); ok && len(val) != 0 {
        man.InstanceId = val
    } else {
        var id *uuid.UUID
        id, err = uuid.NewV4()
        if err != nil {
            return
        }
        man.InstanceId = id.String()
    }

    if val, ok := getInt(m, "API"); ok && val >= 0 {
        man.ApiVersion = val
    } else {
        err = errors.New("API version not found or negative")
        return
    }

    if val, ok := getString(m, "description"); ok && len(val) != 0 {
        man.Description = val
    } else {
        err = errors.New("Description not found")
        return
    }

    if val, ok := getString(m, "version"); ok && len(val) != 0 {
        var ver *semver.Version
        ver, err = semver.NewVersion(val)
        if err != nil {
            err = errors.New("Version number is invalid")
            return
        }
        man.Version = *ver
    } else {
        err = errors.New("Version number not supplied")
        return
    }

    if capMap, ok := getMap(m, "capabilities"); ok {
        var caps []*app.Capability
        caps, err = parseCapabilityMap(capMap)
        if err != nil {
            return
        }
        man.Capabilities = caps
    } else {
        man.Capabilities = make([]*app.Capability, 0)
    }

    if depMap, ok := getMap(m, "dependencies"); ok {
        var deps []*app.Capability
        deps, err = parseCapabilityMap(depMap)
        if err != nil {
            return
        }
        man.Dependencies = deps
    } else {
        man.Dependencies = make([]*app.Capability, 0)
    }
    return
}


func parseCapabilityMap(m jsonData) ([]*app.Capability, error) {
    ret := make([]*app.Capability, 0)
    for k, v := range m {
        switch vv := v.(type) {
        case string:
            ver, err := semver.NewVersion(vv)
            if err != nil {
                return nil, fmt.Errorf("Capability %s has invalid version", k)
            }
            cpb := new(app.Capability)
            cpb.Version = *ver
            cpb.Name = k
            ret = append(ret, cpb)
        default:
            return nil, fmt.Errorf("Capability %s has invalid version", k)
        }
    }
    return ret, nil
}
