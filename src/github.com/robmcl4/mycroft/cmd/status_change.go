package cmd

import (
    "github.com/robmcl4/mycroft/app"
    "github.com/robmcl4/mycroft/registry"
    "log"
    "errors"
    "encoding/json"
)


type StatusChange struct {
    App *app.App
    NewStatus int
    Priority int
}


func NewStatusChange(a *app.App, status int, data []byte) (*Command, error) {
    sc := new(StatusChange)
    sc.App = a
    sc.NewStatus = status
    sc.Priority = -1
    if data != nil {
        var parsed interface{}
        err := json.Unmarshal(data, &parsed)
        if err != nil {
            return nil, err
        }
        m := parsed.(map[string]interface{})
        if val, ok := getInt(m, "priority"); ok {
            sc.Priority = val
        } else {
            return nil, errors.New("Priority was missing or not a number")
        }
    }
    ret := new(Command)
    ret.Execute = sc.Execute
    return ret, nil
}


// change the app's status and notify all those that depend on this app
func (sc *StatusChange) Execute() {
    if sc.App.Status != sc.NewStatus {
        sc.App.Status = sc.NewStatus
        sc.App.Priority = sc.Priority
        if sc.App.Manifest != nil {
            log.Printf("Changing status of %s to '%s'\n", sc.App.Manifest.InstanceId, sc.App.StatusString())
        }
        if sc.App.Manifest != nil {
            // notify everyone who depends on us
            for _, cpb := range sc.App.Manifest.Capabilities {
                for _, dependent := range registry.GetDependents(cpb) {
                    body := make(map[string]interface{})
                    inner := make(map[string]string)
                    inner[sc.App.Manifest.InstanceId] = sc.App.StatusString()
                    body[cpb.Name] = inner
                    dependent.Send("APP_DEPENDENCY", body)
                }
            }
        }
    }
}
