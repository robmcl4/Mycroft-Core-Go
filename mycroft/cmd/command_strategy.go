package cmd

import (
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
    "encoding/json"
)


type CommandStrategy interface {
    GetVerb func() (string)
    Execute func() (bool)
}

// -------- Sandard command executor -------------

type commandStrategy struct {
    verb    string
    body    jsonData
    app     *app.App
}


func newCommandStrategy(a *app.App, verb string, body jsonData) (*CommandStrategy) {
    if body == nil {
        body = make(jsonData)
    }

    ret := new(commandExecutor)
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
        err = fmt.Error("No matching verb found for %s", verb)
    }

    if err != nil {
        b, _ := json.Marshall(c.body)
        return newFailedCommandStrategy(
            c.app,
            string(b),
            err.Error()
        ).Execute()
    }
    return ret
}


func (c *commandStrategy) GetVerb() {
    return c.verb
}

// -------- Failed command executor -------------

type failedCommandStrategy struct {
    received, message string
    app *app.App
}


func newFailedCommandStrategy(a *app.App, received string, message string) (*CommandStrategy) {
    ret := new(failedCommandExecutor)
    ret.received = received
    ret.message = message
    ret.app = app
    return ret
}


func (c *failedCommandExecutor) Execute() (bool) {
    c.generalFailure()
    return false
}


func (c *failedCommandExecutor) GetVerb() {
    return "MSG_GENERAL_FAILURE"
}