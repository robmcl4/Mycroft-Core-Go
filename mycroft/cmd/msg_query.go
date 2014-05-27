package cmd

import (
    "log"
    "encoding/json"
    "errors"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/registry"
)


type MsgQuery struct {
    App *app.App
    Id string
    Capability *app.Capability
    Action string
    Data interface{}
    InstanceIds []string
    Priority int
}


func NewMsgQuery(a *app.App, data []byte) (*Command, error) {
    mq := new(MsgQuery)
    mq.App = a

    // Parse the JSON from the manifest
    var parsed interface{}
    err := json.Unmarshal(data, &parsed)
    if err != nil {
        return nil, err
    }
    m := parsed.(map[string]interface{})

    if val, ok := getString(m, "id"); ok {
        mq.Id = val
    } else {
        return nil, errors.New("ID not supplied")
    }

    if val, ok := getString(m, "capability"); ok {
        found := false
        for _, dep := range a.Manifest.Dependencies {
            if dep.Name == val {
                found = true
                mq.Capability = dep
                break
            }
        }
        if !found {
            return nil, errors.New("This capability was not listed as a dependency")
        }
    } else {
        return nil, errors.New("No capability was given")
    }

    if val, ok := getString(m, "action"); ok {
        mq.Action = val
    } else {
        return nil, errors.New("No action was given")
    }

    mq.Data = m["data"]

    mq.InstanceIds = make([]string, 0)
    if val, ok := m["instanceId"]; ok {
        switch vv := val.(type) {
        case []interface{}:
            for _, inst := range vv {
                switch vinst := inst.(type) {
                case string:
                    mq.InstanceIds = append(mq.InstanceIds, vinst)
                default:
                    return nil, errors.New("InstanceId not a string")
                }
            }
        default:
            return nil, errors.New("instanceId must be an array of strings")
        }
    }

    if val, ok := getInt(m, "priority"); ok {
        mq.Priority = val
    } else {
        return nil, errors.New("Priority was not a valid integer")
    }

    ret := new(Command)
    ret.Execute = mq.Execute
    return ret, nil
}


// send this message query to all targeted apps
func (mq *MsgQuery) Execute() {
    if mq.App.Manifest != nil {
        log.Printf("Processing query from %s\n", mq.App.Manifest.InstanceId)
    }
    registry.RecordMsg(mq.App, mq.Id)
    body := make(map[string]interface{})
    body["fromInstanceId"] = mq.App.Manifest.InstanceId
    body["id"] = mq.Id
    body["priority"] = mq.Priority
    body["data"] = mq.Data
    body["instanceIds"] = mq.InstanceIds
    body["action"] = mq.Action
    body["capability"] = mq.Capability.Name

    // if this is an undirected query
    if len(mq.InstanceIds) == 0 {
        // send to all providers of the capability
        for _, provider := range registry.GetProviders(mq.Capability) {
            provider.Send("MSG_QUERY", body)
        }
    } else {
        // this is a directed query, send to all given instance ids
        for _, instName := range mq.InstanceIds {
            if inst, ok := registry.GetInstance(instName); ok {
                inst.Send("MSG_QUERY", body)
            }
        }
    }
}
