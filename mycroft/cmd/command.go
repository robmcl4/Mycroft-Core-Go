package cmd


import (
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
    "strings"
    "errors"
    "fmt"
    "log"
)


type Command struct {
    Execute func()
}


// Parse a command given the command blob (verb and body, no message length)
func ParseCommand(a *app.App, blob []byte) (*Command) {
    ret, err := internalParseCommand(a, blob)
    if err != nil {
        log.Printf("General failure: %s\n", err.Error())
        recieved := string(blob)
        gf := NewGeneralFailure(a, recieved, err.Error())
        ret = new(Command)
        ret.Execute = gf.Execute
    }
    return ret
}


func internalParseCommand(a *app.App, blob []byte) (*Command, error) {
    str := string(blob)
    first_space := strings.Index(str, " ")
    var verb string
    var body []byte
    if first_space < 0 {
        verb = str
    } else {
        verb = str[:first_space]
        if first_space == len(str)-1 {
            return nil, errors.New("Cannot end verb with space")
        }
        body = []byte(str[first_space+1:])
    }

    switch verb {
    case "APP_MANIFEST":
        return NewAppManifest(a, body)
    case "APP_UP":
        return NewStatusChange(a, app.STATUS_UP, nil)
    case "APP_DOWN":
        return NewStatusChange(a, app.STATUS_DOWN, nil)
    case "APP_IN_USE":
        return NewStatusChange(a, app.STATUS_IN_USE, body)
    case "MSG_QUERY":
        return NewMsgQuery(a, body)
    case "MSG_BROADCAST":
        return NewMsgBroadcast(a, body)
    case "MSG_QUERY_SUCCESS":
        return NewMsgQuerySuccess(a, body)
    case "MSG_QUERY_FAIL":
        return NewMsgQueryFail(a, body)
    }
    return nil, fmt.Errorf("No matching verb found for %s", verb)
}


// get a string from the given map
func getString(m map[string]interface{}, key string) (string, bool) {
    if val, ok := m[key]; ok {
        switch vv := val.(type) {
        case string:
            return vv, true
        default:
            return "", false
        }
    } else {
        return "", false
    }
}


// get an integer from the given map
func getInt(m map[string]interface{}, key string) (int, bool) {
    if val, ok := m[key]; ok {
        switch vv := val.(type) {
        case float64:
            return int(vv), true
        default:
            return 0, false
        }
    } else {
        return 0, false
    }
}


// get a map from the given map
func getMap(m map[string]interface{}, key string) (map[string]interface{}, bool) {
    if val, ok := m[key]; ok {
        switch vv := val.(type) {
        case map[string]interface{}:
            return vv, true
        default:
            return nil, false
        }
    } else {
        return nil, false
    }
}