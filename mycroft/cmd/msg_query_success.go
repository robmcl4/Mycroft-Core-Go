package cmd

import (
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/registry/msg_archive"
    "log"
    "errors"
    "encoding/json"
)


func (c *commandStrategy) msgQuerySuccess() (error) {
    log.Printf("Replying to message from app %s\n", c.app.Manifest.InstanceId)

    var id string
    if id, ok := getString(c.body); !ok {
        return errors.New("No id found")
    }

    ret := c.body["ret"]

    body := make(jsonData)
    body["fromInstanceId"] = mqs.app.Manifest.InstanceId
    body["id"] = id
    body["ret"] = ret
    if recipient, ok := msg_archive.GetMsg(id); ok {
        recipient.Send("MSG_QUERY_SUCCESS", body)
    } else {
        log.Printf("Warning: no app found to reply to for query id %s\n", id)
    }
}
