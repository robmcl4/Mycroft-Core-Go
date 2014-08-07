package msg_archive

import (
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
)

var msgIdMap map[string] *app.App = make(map[string]*app.App)
var appsMessages map[*app.App] []string = make(map[*app.App][]string)


// record that an app sent a message
func RecordMsg(a *app.App, id string) {
    msgIdMap[id] = a
    if val, ok := appsMessages[a]; ok {
        appsMessages[a] = append(val, id)
    } else {
        appsMessages[a] = []string{id}
    }
}


// get a reference to an app that sent a message
func GetMsg(id string) (*app.App, bool) {
    ret, ok := msgIdMap[id]
    return ret, ok
}


// remove all of an app's messages from the archive
func RemoveAppsMessages(a *app.App) {
    if ids, ok := appsMessages[a]; ok {
        for _, id := range ids {
            delete(msgIdMap, id)
        }
    }
    delete(appsMessages, a)
}
