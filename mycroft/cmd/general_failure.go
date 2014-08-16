package cmd

import (
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
)

func generalFailure(a *app.App, recieved string, message string) (error) {
    body := make(jsonData)
    body["recieved"] = recieved
    body["message"] = message
    a.Send("MSG_GENERAL_FAILURE", body)
}
