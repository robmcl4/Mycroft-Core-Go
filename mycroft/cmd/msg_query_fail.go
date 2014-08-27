package cmd

import (
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/registry/msg_archive"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/logging"
    "errors"
)


func (c *commandStrategy) msgQueryFail() (error) {
    c.app.RWMutex.RLock()
    defer c.app.RWMutex.RUnlock()

    logging.Debug("Sending message query fail from %s", c.app.Manifest.InstanceId)
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
        logging.Warning("unrecognized message query id %s from app %s",
                        id,
                        c.app.Manifest.InstanceId)
    }
    return nil
}
