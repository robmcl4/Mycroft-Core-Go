package cmd

import (
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/registry/msg_archive"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/logging"
    "errors"
)


func (c *commandStrategy) msgQuerySuccess() (error) {
    c.app.RWMutex.RLock()
    defer c.app.RWMutex.RUnlock()

    logging.Debug("Replying to message from app %s", c.app.Manifest.InstanceId)

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
        logging.Warning("unrecognized message query id %s from app %s",
                        id,
                        c.app.Manifest.InstanceId)
    }
    return nil
}
