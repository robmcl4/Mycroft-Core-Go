package cmd


import (
    "github.com/robmcl4/mycroft/app"
)


type GeneralFailure struct {
    App *app.App
    Recieved string
    Message string
}


func NewGeneralFailure(a *app.App, recieved string, message string) (*Command) {
    gf := new(GeneralFailure)
    gf.App = a
    gf.Recieved = recieved
    gf.Message = message
    ret := new(Command)
    ret.Execute = gf.Execute
    return ret
}


func (gf *GeneralFailure) Execute() {
    body := make(map[string]interface{})
    body["recieved"] = gf.Recieved
    body["message"] = gf.Message
    gf.App.Send("MSG_GENERAL_FAILURE", body)
}
