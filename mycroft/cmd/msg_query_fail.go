package cmd

import (
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/registry/msg_archive"
    "errors"
    "log"
)


func (c *commandStrategy) msgQueryFail() (error) {
    log.Printf("Sending message query fail from %s", c.app.Manifest.InstanceId)
    var id, message string

    if id_, ok := getString(c.body, "id"); ok {
        id = id_
    } else {
        return errors.New("No id found")
    }

    if message_, ok := getString(c.body, "message"); !ok {
        message = message_
    } else {
        return errors.New("No message found")
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
    return nil
}
