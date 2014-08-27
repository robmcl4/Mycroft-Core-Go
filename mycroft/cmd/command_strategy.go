package cmd

import (
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/logging"
    "encoding/json"
    "fmt"
)


type CommandStrategy interface {
    GetVerb () string
    Execute () bool
}

// -------- Sandard command executor -------------

type commandStrategy struct {
    verb    string
    body    jsonData
    app     *app.App
}


func newCommandStrategy(a *app.App, verb string, body jsonData) (CommandStrategy) {
    if body == nil {
        body = make(jsonData)
    }

    ret := new(commandStrategy)
    ret.verb = verb
    ret.body = body
    ret.app = a
    return ret
}


func (c *commandStrategy) Execute() (bool) {
    var err error
    ret := true

    switch c.verb {
    case "APP_MANIFEST":
        ret = c.appManifest()
    case "APP_UP", "APP_DOWN", "APP_IN_USE":
        err = c.statusChange()
    case "MSG_QUERY":
        err = c.msgQuery()
    case "MSG_BROADCAST":
        err = c.msgBroadcast()
    case "MSG_QUERY_SUCCESS":
        err = c.msgQuerySuccess()
    case "MSG_QUERY_FAIL":
        err = c.msgQueryFail()
    default:
        err = fmt.Errorf("No matching verb found for %s", c.verb)
    }

    if err != nil {
        id := "NO_ID_FOUND"
        if c.app.Manifest != nil {
            id = c.app.Manifest.InstanceId
        }
        logging.Error("Application %s had command error %s", id, err.Error())
        b, _ := json.Marshal(c.body)
        return newFailedCommandStrategy(c.app, string(b), err.Error()).Execute()
    }
    return ret
}


func (c *commandStrategy) GetVerb() string {
    return c.verb
}

// -------- Failed command executor -------------

type failedCommandStrategy struct {
    received string
    message string
    app *app.App
}


func newFailedCommandStrategy(a *app.App, received string, message string) (CommandStrategy) {
    ret := new(failedCommandStrategy)
    ret.received = received
    ret.message = message
    ret.app = a
    return ret
}


func (c *failedCommandStrategy) Execute() (bool) {
    c.generalFailure()
    return false
}


func (c *failedCommandStrategy) GetVerb() string {
    return "MSG_GENERAL_FAILURE"
}
