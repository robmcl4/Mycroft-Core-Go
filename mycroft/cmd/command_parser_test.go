package cmd

import (
    "github.com/stretchr/testify/assert"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
    "testing"
)


func TestJsonDataType(t *testing.T) {
    m := make(jsonData)
    m["foo"] = "bar"
    assert.Equal(t, m["foo"], "bar")
}


func TestParseBodyWithValidData(t *testing.T) {
    var got jsonData
    var err error

    got, err = parseBody("{}")
    assert.Nil(t, err)

    got, err = parseBody("{\"foo\":\"bar\"}")
    assert.Nil(t, err)
    assert.Equal(t, got["foo"], "bar")
}


func TestParseBodyWithInvalidData(t *testing.T) {
    var err error

    _, err = parseBody("foobar")
    assert.Error(t, err)

    _, err = parseBody("")
    assert.Error(t, err)

    _, err = parseBody("[\"a\"]")
    assert.Error(t, err)
}


func TestSeparateVerb(t *testing.T) {
    var verb, body string

    verb, body = separateVerb("foo")
    assert.Equal(t, verb, "foo")
    assert.Equal(t, body, "")

    verb, body = separateVerb("foo bar")
    assert.Equal(t, verb, "foo")
    assert.Equal(t, body, "bar")

    verb, body = separateVerb("")
    assert.Equal(t, verb, "")
    assert.Equal(t, body, "")
}


func TestParseVerbAndBody(t *testing.T) {
    var verb string
    var body jsonData
    var err error

    verb, body, err = parseVerbAndBody("FOO {\"bar\":\"baz\"}")
    assert.NoError(t, err)
    assert.Equal(t, verb, "FOO")
    assert.NotNil(t, body)
    assert.Equal(t, body["bar"], "baz")

    verb, body, err = parseVerbAndBody("SLAM")
    assert.NoError(t, err)
    assert.Equal(t, verb, "SLAM")
    assert.Nil(t, body)

    verb, body, err = parseVerbAndBody("FOO []")
    assert.Error(t, err)

    verb, body, err = parseVerbAndBody("")
    assert.Error(t, err)
}


func TestInternalParseCommand(t *testing.T) {
    var strat CommandStrategy
    var err error
    var fakeApp *app.App

    strat, err = internalParseCommand(fakeApp, "FOO {\"bar\":\"baz\"}")
    assert.NoError(t, err)
    assert.NotNil(t, strat)
    assert.IsType(t, &commandStrategy{}, strat)

    strat, err = internalParseCommand(fakeApp, "FOO")
    assert.NoError(t, err)
    assert.NotNil(t, strat)
    assert.IsType(t, &commandStrategy{}, strat)

    strat, err = internalParseCommand(fakeApp, "invalid data []")
    assert.Error(t, err)

    strat, err = internalParseCommand(fakeApp, "")
    assert.Error(t, err)
}


func TestParseCommand(t *testing.T) {
    var strat CommandStrategy
    var status bool
    fakeApp := &app.App{}

    strat, status = ParseCommand(fakeApp, "FOO {\"bar\":\"baz\"}")
    assert.True(t, status)
    assert.NotNil(t, strat)
    assert.IsType(t, &commandStrategy{}, strat)


    strat, status = ParseCommand(fakeApp, "FOO")
    assert.True(t, status)
    assert.NotNil(t, strat)
    assert.IsType(t, &commandStrategy{}, strat)


    strat, status = ParseCommand(fakeApp, "")
    assert.False(t, status)
    assert.IsType(t, &failedCommandStrategy{}, strat)
}


func TestGetString(t *testing.T) {
    var toPut jsonData
    var got string
    var ok bool

    toPut, _ = parseBody("{\"foo\":\"bar\"}")
    got, ok = getString(toPut, "foo")
    assert.Equal(t, got, "bar")
    assert.True(t, ok)

    toPut, _ = parseBody("{\"foo\":12}")
    got, ok = getString(toPut, "foo")
    assert.False(t, ok)

    toPut, _ = parseBody("{}")
    got, ok = getString(toPut, "foo")
    assert.False(t, ok)
}


func TestGetInt(t *testing.T) {
    var toPut jsonData
    var got int
    var ok bool

    toPut, _ = parseBody("{\"foo\":12}")
    got, ok = getInt(toPut, "foo")
    assert.True(t, ok)
    assert.Equal(t, got, 12)

    toPut, _ = parseBody("{\"foo\":[]}")
    got, ok = getInt(toPut, "foo")
    assert.False(t, ok)

    toPut, _ = parseBody("{}")
    got, ok = getInt(toPut, "foo")
    assert.False(t, ok)
}


func TestGetMap(t *testing.T) {
    var toPut jsonData
    var got jsonData
    var ok bool

    toPut, _ = parseBody("{\"foo\":{\"bar\":\"baz\"}}")
    got, ok = getMap(toPut, "foo")
    assert.True(t, ok)
    assert.Equal(t, got["bar"], "baz")

    toPut, _ = parseBody("{\"foo\":12}")
    got, ok = getMap(toPut, "foo")
    assert.False(t, ok)

    toPut, _ = parseBody("{}")
    got, ok = getMap(toPut, "foo")
    assert.False(t, ok)
}
