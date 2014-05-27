package cmd

import (
    "github.com/robmcl4/mycroft/app"
    "github.com/robmcl4/mycroft/registry"
    "errors"
    "log"
    "encoding/json"
)


type MsgQueryFail struct {
    App *app.App
    Id string
    Message string
}


func NewMsgQueryFail(a *app.App, data []byte) (*Command, error) {
    mqf := new(MsgQueryFail)
    mqf.App = a

    // Parse the JSON from the manifest
    var parsed interface{}
    err := json.Unmarshal(data, &parsed)
    if err != nil {
        return nil, err
    }
    m := parsed.(map[string]interface{})

    if val, ok := getString(m, "id"); ok {
        mqf.Id = val
    } else {
        return nil, errors.New("No id found")
    }

    if val, ok := getString(m, "message"); ok {
        mqf.Message = val
    } else {
        return nil, errors.New("No message found")
    }

    ret := new(Command)
    ret.Execute = mqf.Execute
    return ret, nil
}


func (mqf *MsgQueryFail) Execute() {
    log.Printf("Sending message query fail from %s", mqf.App.Manifest.InstanceId)

    body := make(map[string]interface{})
    body["fromInstanceId"] = mqf.App.Manifest.InstanceId
    body["id"] = mqf.Id
    body["message"] = mqf.Message
    if recipient, ok := registry.GetMsg(mqf.Id); ok {
        recipient.Send("MSG_QUERY_FAIL", body)
    } else {
        log.Printf("Warning: no app found to reply to for query id %s\n", mqf.Id)
    }
}
