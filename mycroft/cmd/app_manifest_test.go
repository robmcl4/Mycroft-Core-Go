package cmd

import (
    "github.com/stretchr/testify/assert"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
    "testing"
)


func TestParseCapabilityMap(t *testing.T) {
    var toParse jsonData
    var got []*app.Capability
    var err error

    toParse, _ = parseBody("{}")
    got, err = parseCapabilityMap(toParse)
    assert.NotNil(t, got)
    assert.NoError(t, err)
    assert.Len(t, got, 0)

    toParse, _ = parseBody(`{"foo":"1.1.1"}`)
    got, err = parseCapabilityMap(toParse)
    assert.NotNil(t, got)
    assert.NoError(t, err)
    assert.Len(t, got, 1)
    assert.Equal(t, got[0].Name, "foo")
    assert.Equal(t, got[0].Version.Major, 1)
    assert.Equal(t, got[0].Version.Minor, 1)
    assert.Equal(t, got[0].Version.Patch, 1)
}


func TestDecodeManifest(t *testing.T) {
    base := `{
        "name": "foo_name",
        "displayName": "foo_display_name",
        "instanceId": "inst 007",
        "API": 1,
        "description": "a description",
        "version": "1.1.2",
        "capabilities": {
            "foo": "1.1.3"
        },
        "dependencies": {
            "bar": "1.2.3"
        }
    }`

    toParse, err := parseBody(base)
    assert.NoError(t, err)

    got, err := decodeManifest(toParse)
    assert.NoError(t, err)
    assert.NotNil(t, got)

    assert.Equal(t, got.Name, "foo_name")
    assert.Equal(t, got.DisplayName, "foo_display_name")
    assert.Equal(t, got.InstanceId, "inst 007")
    assert.Equal(t, got.ApiVersion, 1)
    assert.Equal(t, got.Description, "a description")
    assert.Equal(t, got.Version.Major, 1)
    assert.Equal(t, got.Version.Minor, 1)
    assert.Equal(t, got.Version.Patch, 2)
    assert.Len(t, got.Capabilities, 1)
    assert.Len(t, got.Dependencies, 1)

    // make sure it complains without name
    toParse, err = parseBody(base)
    delete(toParse, "name")
    got, err = decodeManifest(toParse)
    assert.Error(t, err)

    // make sure it complains without displayName
    toParse, err = parseBody(base)
    delete(toParse, "displayName")
    got, err = decodeManifest(toParse)
    assert.Error(t, err)

    // make sure an empty instanceId is filled in
    toParse, err = parseBody(base)
    delete(toParse, "instanceId")
    got, err = decodeManifest(toParse)
    assert.NotNil(t, got.InstanceId)
}
