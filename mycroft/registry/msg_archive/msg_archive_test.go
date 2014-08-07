package msg_archive

import (
    "testing"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
    "github.com/stretchr/testify/assert"
)

func setUp() {
    msgIdMap = make(map[string]*app.App)
    appsMessages = make(map[*app.App][]string)
}

func TestRecordMsgAddsToMap(t *testing.T) {
    setUp()
    a := app.NewApp()
    id := "something"
    RecordMsg(a, id)

    assert.Len(t, msgIdMap, 1, "The map should have 1 id")
    assert.Equal(t, msgIdMap[id], a, "The map should have some_id -> a")
}

func TestRecordMsgCreatesAppsMessagesKey(t *testing.T) {
    setUp()
    a := app.NewApp()
    id := "something"
    RecordMsg(a, id)

    assert.Len(t, appsMessages, 1, "There should be 1 app registered")
    assert.Equal(t, appsMessages[a], []string{id})
}

func TestRecordMsgAppendsAppsMessages(t *testing.T) {
    setUp()
    a := app.NewApp()
    id := "something"
    appsMessages[a] = []string{"something_else"}
    RecordMsg(a, id)

    assert.Len(t, appsMessages[a], 2)
}

func TestRemoveAppsMessagesWithNothingRegistered(t *testing.T) {
    setUp()
    a := app.NewApp()
    assert.NotPanics(t, func() { RemoveAppsMessages(a) })
}

func TestRemoveAppsMessagesWithMessagesRegistered(t *testing.T) {
    setUp()
    a := app.NewApp()
    RecordMsg(a, "foo")
    RecordMsg(a, "bar")
    RemoveAppsMessages(a)
    assert.Len(t, appsMessages, 0)
    assert.Len(t, msgIdMap, 0)
}

func TestRemoveAppsMessagesWithMessagesFromOthersRegistered(t *testing.T) {
    setUp()
    a1 := app.NewApp()
    a2 := app.NewApp()
    RecordMsg(a1, "foo")
    RecordMsg(a1, "bar")
    RecordMsg(a2, "baz")
    RecordMsg(a2, "qux")
    RemoveAppsMessages(a1)
    assert.Len(t, appsMessages[a2], 2)
    assert.Equal(t, appsMessages[a2], []string{"baz", "qux"})
    assert.Len(t, msgIdMap, 2)
    assert.Equal(t, msgIdMap["baz"], a2)
}
