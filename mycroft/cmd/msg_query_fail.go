package cmd

import (
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/registry/msg_archive"
    "errors"
    "log"
    "encoding/json"
)


func (c *commandStrategy) msgQueryFail() (error) {
    log.Printf("Sending message query fail from %s", mqf.App.Manifest.InstanceId)
    var id, message string

    if id, ok := getString(c.body, "id"); !ok {
        return nil, errors.New("No id found")
    }
    if message, ok := getString(c.body, "message"); !ok {
        return nil, errors.New("No message found")
    }

    body := make(jsonData)
    body["fromInstanceId"] = c.app.Manifest.InstanceId
    body["id"] = id
    body["message"] = message
    if recipient, ok := msg_archive.GetMsg(id); ok {
        recipient.Send("MSG_QUERY_FAIL", body)
    } else {
        log.Printf("Warning: unrecognized message id %s\n", id)
    }
}
