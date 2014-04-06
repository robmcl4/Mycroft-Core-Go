package cmd


import (
    "github.com/robmcl4/mycroft/app"
    "github.com/robmcl4/mycroft/registry"
    "github.com/coreos/go-semver/semver"
    "github.com/nu7hatch/gouuid"
    "errors"
    "log"
    "fmt"
    "encoding/json"
)

type AppManifest struct {
    App *app.App
    data []byte
}


func NewAppManifest(a *app.App, data []byte) (*Command, error) {
    am := new(AppManifest)
    am.App = a
    am.data = data
    ret := new(Command)
    ret.Execute = am.Execute
    return ret, nil
}


// construct the app based on its manifest
func (a *AppManifest) Execute() {
    log.Println("Parsing application's manifest")
    man, err := decodeManifest(a.data)
    if err != nil {
        log.Printf("ERROR: app's manifest is invalid: %s", err)
        go sendManifestFail(a.App, err.Error())
        return
    }
    if _, exists := registry.GetInstance(man.InstanceId); exists {
        log.Printf("ERROR: app's instance id is already taken: %s", man.InstanceId)
        go sendManifestFail(a.App, fmt.Sprintf("instanceId %s is in use", man.InstanceId))
        return
    }
    a.App.Status = app.STATUS_DOWN
    a.App.Manifest = man
    registry.Register(a.App)
    sendManifestOkAndDependencies(a.App)
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


func decodeManifest(data []byte) (man *app.Manifest, err error) {
    man = new(app.Manifest)

    // Parse the JSON from the manifest
    var parsed interface{}
    err = json.Unmarshal(data, &parsed)
    if err != nil {
        return
    }
    m := parsed.(map[string]interface{})

    // start loading in values from the manifest
    if str, ok := getString(m, "name"); ok && len(str) != 0 {
        man.Name = str
    } else {
        err = errors.New("No name was found")
        return
    }

    if val, ok := getString(m, "name"); ok && len(val) != 0 {
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
        man.Version = ver
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


func parseCapabilityMap(m map[string]interface{}) ([]*app.Capability, error) {
    ret := make([]*app.Capability, 0)
    for k, v := range m {
        switch vv := v.(type) {
        case string:
            ver, err := semver.NewVersion(vv)
            if err != nil {
                return nil, fmt.Errorf("Capability %s has invalid version", k)
            }
            cpb := new(app.Capability)
            cpb.Version = ver
            cpb.Name = k
            ret = append(ret, cpb)
        default:
            return nil, fmt.Errorf("Capability %s has invalid version", k)
        }
    }
    return ret, nil
}
