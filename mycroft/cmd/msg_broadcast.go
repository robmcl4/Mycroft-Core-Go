package cmd

import (
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/registry"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/registry/msg_archive"
    "log"
    "errors"
)

func msgBroadcast(a *app.App, body jsonData) (error) {
    var id string
    content := body["content"]

    log.Printf("Sending message broadcast from %s\n", a.Manifest.InstanceId)

    if id, ok := getString(m, "id"); !ok {
        return errors.New("No id found")
    }

    msg_archive.RecordMsg(a, mb.Id)
    toSend := make(jsonData)
    toSend["fromInstanceId"] = a.Manifest.InstanceId
    toSend["id"] = id
    toSend["content"] = content
    for _, cpb := range a.Manifest.Capabilities {
        for _, dep := range registry.GetDependents(cpb) {
            dep.Send("MSG_BROADCAST", toSend)
        }
    }
}
