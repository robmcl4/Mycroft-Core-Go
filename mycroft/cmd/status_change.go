package cmd

import (
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/registry"
    "log"
    "errors"
)


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
            for _, cpb := range c.app.Manifest.Capabilities {
                for _, dependent := range registry.GetDependents(cpb) {
                    body := make(map[string]interface{})
                    inner := make(map[string]string)
                    inner[c.app.Manifest.InstanceId] = c.app.StatusString()
                    body[cpb.Name] = inner
                    dependent.Send("APP_DEPENDENCY", body)
                }
            }
        }
    }
    return nil
}
