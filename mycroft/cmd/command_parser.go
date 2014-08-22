// Package cmd parses application commands and defines command logic
package cmd

import (
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
    "strings"
    "errors"
    "fmt"
    "log"
    "encoding/json"
)


type jsonData map[string]interface{}


// Parses a command given the command's message string
// (verb and body, no message length)
func ParseCommand(a *app.App, message string) (*Strategy, bool) {
    ret, err := internalParseCommand(a, message)
    if err != nil {
        log.Printf("General failure: %s\n", err.Error())
        return newFailedCommandStrategy(a, message, err.Error()), false
    }
    return ret, true
}


func internalParseCommand(a *app.App, message string) (*Strategy, error) {
    verb, body, err := parseVerbAndBody(message)
    if err != nil {
        return nil, err
    }

    return newCommandStrategy(app, verb, body)
}


// Parses the command message into the verb and body. The body may be nil.
func parseVerbAndBody(message string) (verb string, body jsonData, err error) {
    verb, maybeBody := separateVerb(message)
    if maybeBody != nil {
        body, err = parseBody(maybeBody)
    }
}


// Separates the message verb from the body
func separateVerb(message string) (verb string, body string) {
    split := strings.SplitN(message, " ", 2)
    verb = split[0]
    if len(split) == 2 {
        body = split[1]
    }
}


// Parses the message body to a JSON map
func parseBody(body string) (jsonData, error) {
    var parsed interface{}
    err := json.Unmarshal(data, &parsed)
    if err != nil {
        return nil, err
    }
    return parsed.(JSON), nil
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