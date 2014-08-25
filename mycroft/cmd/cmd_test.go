package cmd

import (
    "github.com/stretchr/testify/assert"
    "testing"
)


func TestJsonDataType(t *testing.T) {
    m := make(jsonData)
    m["foo"] = "bar"
    assert.Equal(t, m["foo"], "bar")
}


func TestParseBody(t *testing.T) {
    var got jsonData
    var err error

    got, err = parseBody("{}")
    assert.Nil(t, err)

    got, err = parseBody("{\"foo\":\"bar\"}")
    assert.Nil(t, err)
    assert.Equal(t, got["foo"], "bar")
}
