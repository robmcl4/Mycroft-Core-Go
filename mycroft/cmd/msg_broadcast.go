package cmd

import (
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/registry"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/registry/msg_archive"
    "log"
    "errors"
)

func (c *commandStrategy) msgBroadcast() (error) {
    var id string
    content := c.body["content"]

    log.Printf("Sending message broadcast from %s\n", c.app.Manifest.InstanceId)

    if id_, ok := getString(c.body, "id"); ok {
        id = id_
    } else {
        return errors.New("No id found")
    }

    msg_archive.RecordMsg(c.app, id)
    toSend := make(jsonData)
    toSend["fromInstanceId"] = c.app.Manifest.InstanceId
    toSend["id"] = id
    toSend["content"] = content
    for _, cpb := range c.app.Manifest.Capabilities {
        for _, dep := range registry.GetDependents(cpb) {
            dep.Send("MSG_BROADCAST", toSend)
        }
    }
    return nil
}
