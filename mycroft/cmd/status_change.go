package cmd

import (
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/registry"
    "log"
    "errors"
)


// Changes the status of App a to the given status and priority.
// By design priority is ignored for all statuses other than STATUS_IN_USE.
func ChangeAppStatus(a *app.App, stat int, prio int) {
    a.Status = stat
    a.Priority = prio
    sendDependencyNotice(a)
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
            sendDependencyNotice(c.app)
        }
    }
    return nil
}


func sendDependencyNotice(a *app.App) {
    // notify everyone who depends on us
    for _, cpb := range a.Manifest.Capabilities {
        for _, dependent := range registry.GetDependents(cpb) {
            body := make(jsonData)
            inner := make(jsonData)
            inner[a.Manifest.InstanceId] = a.StatusString()
            body[cpb.Name] = inner
            dependent.Send("APP_DEPENDENCY", body)
        }
    }
}
