package cmd

import (
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/registry"
    "log"
    "errors"
    "encoding/json"
)


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
func (c *commandStrategy) statusChange() (error) {
    var newStatus int
    priority := -1

    switch c.verb {
    case "APP_UP":
        newStatus = app.STATUS_UP
    case "APP_DOWN":
        newStatus = app.STATUS_DOWN
    case "APP_IN_USE":
        newStatus = app.STATUS_IN_USE
    default:
        return errors.New("Unrecognized status verb")
    }

    if c.app.Status != newStatus {
        c.app.Status = newStatus
        c.app.Priority = priority

        if c.app.Manifest != nil {
            log.Printf("Changing status of %s to '%s'\n", c.app.Manifest.InstanceId, c.app.StatusString())

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
