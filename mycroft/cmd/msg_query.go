package cmd

import (
    "log"
    "errors"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/registry"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/registry/msg_archive"
)


type parsedMsgQuery struct {
    id string
    capability *app.Capability
    action string
    data interface{}
    instanceIds []string
    priority int
}


func parseMsgQuery(c *commandStrategy) (*parsedMsgQuery, error) {
    ret := new(parsedMsgQuery)

    // get the id
    if id, ok := getString(c.body, "id"); ok {
        return nil, errors.New("ID not supplied")
    } else {
        ret.id = id
    }

    // get the capability, and find the matching app.Capability
    // in the manifest (so we know what version they want)
    if val, ok := getString(c.body, "capability"); ok {
        for _, dep := range c.app.Manifest.Dependencies {
            if dep.Name == val {
                ret.capability = dep
                break
            }
        }
        if ret.capability != nil {
            return nil, errors.New("This capability was not listed as a dependency")
        }
    } else {
        return nil, errors.New("No capability was given")
    }

    // get the action they want to perform
    if action, ok := getString(c.body, "action"); ok {
        ret.action = action
    } else {
        return nil, errors.New("No action was given")
    }

    // get data; can be anything, just grab whatever's there
    ret.data = c.body["data"]

    // get the instance IDs and do a whole bunch of type-checking
    ret.instanceIds = make([]string, 0)
    if val, ok := c.body["instanceId"]; ok {
        switch vv := val.(type) {
        case []interface{}:
            for _, inst := range vv {
                switch vinst := inst.(type) {
                case string:
                    ret.instanceIds = append(ret.instanceIds, vinst)
                default:
                    return nil, errors.New("InstanceId not a string")
                }
            }
        default:
            return nil, errors.New("instanceId must be an array of strings")
        }
    }

    // get the message priority
    if priority, ok := getInt(c.body, "priority"); ok {
        ret.priority = priority
    } else {
        return nil, errors.New("Priority was not a valid integer")
    }

    return ret, nil
}


// send this message query to all targeted apps
func (c *commandStrategy) msgQuery() (error) {
    c.app.RWMutex.RLock()
    defer c.app.RWMutex.RUnlock()

    log.Printf("Processing query from %s\n", c.app.Manifest.InstanceId)
    mq, err := parseMsgQuery(c)
    if err != nil {
        return err
    }

    msg_archive.RecordMsg(c.app, mq.id)
    body := make(jsonData)
    body["fromInstanceId"] = c.app.Manifest.InstanceId
    body["id"] = mq.id
    body["priority"] = mq.priority
    body["data"] = mq.data
    body["instanceId"] = mq.instanceIds
    body["action"] = mq.action
    body["capability"] = mq.capability.Name

    // if this is an undirected query
    if len(mq.instanceIds) == 0 {
        // send to all providers of the capability
        for _, provider := range registry.GetProviders(mq.capability) {
            provider.Send("MSG_QUERY", body)
        }
    } else {
        // this is a directed query, send to all given instance ids
        for _, instName := range mq.instanceIds {
            if inst, ok := registry.GetInstance(instName); ok {
                inst.Send("MSG_QUERY", body)
            }
        }
    }
    return nil
}
