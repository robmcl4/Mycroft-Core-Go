package cmd

import (
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/registry"
    "log"
    "errors"
    "encoding/json"
)


type MsgBroadcast struct {
    App *app.App
    Id string
    Content interface{}
}


func NewMsgBroadcast(a *app.App, data []byte) (*Command, error) {
    mb := new(MsgBroadcast)
    mb.App = a

    // Parse the JSON from the manifest
    var parsed interface{}
    err := json.Unmarshal(data, &parsed)
    if err != nil {
        return nil, err
    }
    m := parsed.(map[string]interface{})

    if val, ok := getString(m, "id"); ok {
        mb.Id = val
    } else {
        return nil, errors.New("No id found")
    }

    mb.Content = m["content"]

    ret := new(Command)
    ret.Execute = mb.Execute
    return ret, nil
}


func (mb *MsgBroadcast) Execute() {
    log.Printf("Sending message broadcast from %s\n", mb.App.Manifest.InstanceId)
    registry.RecordMsg(mb.App, mb.Id)
    body := make(map[string]interface{})
    body["fromInstanceId"] = mb.App.Manifest.InstanceId
    body["id"] = mb.Id
    body["content"] = mb.Content
    for _, cpb := range mb.App.Manifest.Capabilities {
        for _, dep := range registry.GetDependents(cpb) {
            dep.Send("MSG_BROADCAST", body)
        }
    }
}
