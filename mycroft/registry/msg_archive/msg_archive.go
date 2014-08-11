package msg_archive

import (
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
    "sync"
)

var msgIdMap map[string] *app.App = make(map[string]*app.App)
var appsMessages map[*app.App] []string = make(map[*app.App][]string)
var lock *sync.RWMutex = new(sync.RWMutex)

// record that an app sent a message
func RecordMsg(a *app.App, id string) {
    lock.Lock()
    defer lock.Unlock()

    msgIdMap[id] = a
    appsMessages[a] = append(appsMessages[a], id)
}


// get a reference to an app that sent a message
func GetMsg(id string) (*app.App, bool) {
    lock.RLock()
    defer lock.RUnlock()

    ret, ok := msgIdMap[id]
    return ret, ok
}


// remove all of an app's messages from the archive
func RemoveAppsMessages(a *app.App) {
    lock.Lock()
    defer lock.Unlock()

    if ids, ok := appsMessages[a]; ok {
        for _, id := range ids {
            delete(msgIdMap, id)
        }
    }
    delete(appsMessages, a)
}
