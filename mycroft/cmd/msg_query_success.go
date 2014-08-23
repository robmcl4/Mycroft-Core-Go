package cmd

import (
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/registry/msg_archive"
    "log"
    "errors"
)


func (c *commandStrategy) msgQuerySuccess() (error) {
    log.Printf("Replying to message from app %s\n", c.app.Manifest.InstanceId)

    var id string
    if id_, ok := getString(c.body, "id"); ok {
        id = id_
    } else {
        return errors.New("No id found")
    }

    ret := c.body["ret"]

    body := make(jsonData)
    body["fromInstanceId"] = c.app.Manifest.InstanceId
    body["id"] = id
    body["ret"] = ret
    if recipient, ok := msg_archive.GetMsg(id); ok {
        recipient.Send("MSG_QUERY_SUCCESS", body)
    } else {
        log.Printf("Warning: no app found to reply to for query id %s\n", id)
    }
    return nil
}
