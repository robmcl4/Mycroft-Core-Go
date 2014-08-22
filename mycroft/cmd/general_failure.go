package cmd

import (
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
)

func (f *failedCommandStrategy) generalFailure() {
    body := make(jsonData)
    body["recieved"] = recieved
    body["message"] = message
    f.app.Send("MSG_GENERAL_FAILURE", body)
}
