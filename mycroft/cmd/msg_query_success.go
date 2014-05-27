package cmd

import (
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/registry"
    "log"
    "errors"
    "encoding/json"
)


type MsgQuerySuccess struct {
    App *app.App
    Id string
    Ret interface{}
}


func NewMsgQuerySuccess(a *app.App, data []byte) (*Command, error) {
    mqs := new(MsgQuerySuccess)
    mqs.App = a

    // Parse the JSON from the manifest
    var parsed interface{}
    err := json.Unmarshal(data, &parsed)
    if err != nil {
        return nil, err
    }
    m := parsed.(map[string]interface{})

    if val, ok := getString(m, "id"); ok {
        mqs.Id = val
    } else {
        return nil, errors.New("No id found")
    }

    mqs.Ret = m["ret"]

    ret := new(Command)
    ret.Execute = mqs.Execute
    return ret, nil
}


func (mqs *MsgQuerySuccess) Execute() {
    log.Printf("Replying to message from app %s\n", mqs.App.Manifest.InstanceId)
    body := make(map[string]interface{})
    body["fromInstanceId"] = mqs.App.Manifest.InstanceId
    body["id"] = mqs.Id
    body["ret"] = mqs.Ret
    if recipient, ok := registry.GetMsg(mqs.Id); ok {
        recipient.Send("MSG_QUERY_SUCCESS", body)
    } else {
        log.Printf("Warning: no app found to reply to for query id %s\n", mqs.Id)
    }
}
