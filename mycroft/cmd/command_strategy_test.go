package cmd

import (
    "github.com/stretchr/testify/assert"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
    "testing"
)


func TestNewCommandStrategy(t *testing.T) {
    fakeApp := &app.App{}
    body := make(jsonData)

    got := newCommandStrategy(fakeApp, "FOO", body).(*commandStrategy)
    assert.NotNil(t, got)
    assert.Equal(t, got.body, body)
    assert.Equal(t, got.app, fakeApp)
    assert.Equal(t, got.verb, "FOO")

    got = newCommandStrategy(fakeApp, "FOO", nil).(*commandStrategy)
    assert.NotNil(t, got.body)
}


func TestCommandStrategyGetVerb(t *testing.T) {
    toTest := &commandStrategy{}
    toTest.verb = "foobar"
    assert.Equal(t, toTest.GetVerb(), "foobar")
}


func TestNewFailedCommandStrategy(t *testing.T) {
    fakeApp := &app.App{}

    got := newFailedCommandStrategy(fakeApp, "received", "message").(*failedCommandStrategy)
    assert.NotNil(t, got)
    assert.Equal(t, got.app, fakeApp)
    assert.Equal(t, got.received, "received")
    assert.Equal(t, got.message, "message")
}


func TestFailedCommandStrategyGetVerb(t *testing.T) {
    toTest := &failedCommandStrategy{}
    assert.Equal(t, toTest.GetVerb(), "MSG_GENERAL_FAILURE")
}
