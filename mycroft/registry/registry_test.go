package registry


import (
    "testing"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
    "github.com/coreos/go-semver/semver"
    "github.com/stretchr/testify/assert"
)


func getMockApp() (*app.App) {
    // construct a new app
    a := new(app.App)
    a.Manifest = new(app.Manifest)
    a.Manifest.Name = "name1"
    a.Manifest.InstanceId = "instance1"

    // construct capabilities
    caps := make([]*app.Capability, 2)
    cap1 := new(app.Capability)
    v1, _ := semver.NewVersion("1.1.1")
    v2, _ := semver.NewVersion("1.1.2")
    cap1.Version = *v1
    cap1.Name = "capability1"
    caps[0] = cap1
    cap2 := new(app.Capability)
    cap2.Version = *v2
    cap2.Name = "capability2"
    caps[1] = cap2
    a.Manifest.Capabilities = caps

    // construct dependencies
    deps := make([]*app.Capability, 2)
    dep1 := new(app.Capability)
    dep1.Version = *v1
    dep1.Name = "dependency1"
    deps[0] = dep1
    dep2 := new(app.Capability)
    dep2.Version = *v2
    dep2.Name = "dependency2"
    deps[1] = dep2

    a.Manifest.Dependencies = deps

    return a
}


func setUp() {
    capabilitySuppliers = make(map[app.Capability][]*app.App)
    capabilityDependents = make(map[app.Capability][]*app.App)
    instances = make(map[string]*app.App)
}


func TestCapabilitySuppliersMap(t *testing.T) {
    setUp()
    a := getMockApp()
    Register(a)
    // make sure the registry entries exist
    assert.Len(t, capabilitySuppliers, 2, "Two capabilities should be registered")
    for _, capability := range a.Manifest.Capabilities {
        assert.Equal(t, capabilitySuppliers[*capability], []*app.App{a})
    }
}

func TestAppendCapabilitySuppliersMap(t *testing.T) {
    setUp()
    a1 := getMockApp()
    a2 := getMockApp()
    Register(a1)
    Register(a2)
    for _, capability := range a1.Manifest.Capabilities {
        assert.Len(t, capabilitySuppliers[*capability], 2, "Two apps should have this capability")
    }
}


func TestCapabilityDependentsMap(t *testing.T) {
    setUp()
    a := getMockApp()
    Register(a)
    // make sure the registry entries exist
    assert.Len(t, capabilityDependents, 2, "Two dependencies should be registered")
    for _, dependency := range a.Manifest.Dependencies {
        assert.Equal(t, capabilityDependents[*dependency], []*app.App{a})
    }
}


func TestAppendCapabilityDependentsMap(t *testing.T) {
    setUp()
    a1 := getMockApp()
    a2 := getMockApp()
    Register(a1)
    Register(a2)
    for _, dependency := range a1.Manifest.Dependencies {
        assert.Len(t, capabilityDependents[*dependency], 2, "Two apps should have this dependency")
    }
}


func TestInstanceIdMap(t *testing.T) {
    setUp()
    a := getMockApp()
    Register(a)
    assert.Equal(t, instances[a.Manifest.InstanceId], a)
}


func TestGetProviders(t *testing.T) {
    setUp()
    a := getMockApp()
    Register(a)
    apps := GetProviders(a.Manifest.Capabilities[0])
    assert.Len(t, apps, 1, "one app should have this capability")
    assert.Equal(t, apps[0], a, "it should be the correct app")
    // make sure it actually returns a _copy_ of the array
    apps[0] = nil
    apps = GetProviders(a.Manifest.Capabilities[0])
    assert.Equal(t, apps[0], a, "GetProviders should return a copy of the array")
}


func TestGetDependents(t *testing.T) {
    setUp()
    a := getMockApp()
    Register(a)
    apps := GetDependents(a.Manifest.Dependencies[0])
    assert.Len(t, apps, 1)
    assert.Equal(t, apps[0], a)
    // make sure it actually returns a _copy_ of the array
    apps[0] = nil
    apps = GetDependents(a.Manifest.Dependencies[0])
    assert.Equal(t, apps[0], a, "GetDependents should return a copy of the array")
}


func TestGetInstance(t *testing.T) {
    setUp()
    a := getMockApp()
    Register(a)
    inst, ok := GetInstance(a.Manifest.InstanceId)
    assert.True(t, ok)
    assert.Equal(t, inst, a)
    inst, ok = GetInstance("foo")
    assert.False(t, ok, "should not have found an instance")
}


func TestRemoveApp(t *testing.T) {
    setUp()
    a := getMockApp()
    b := getMockApp()
    b.Manifest.InstanceId = "foo11223"
    Register(a)
    Register(b)
    Remove(a)
    inst, ok := GetInstance(b.Manifest.InstanceId)
    assert.True(t, ok)
    assert.Equal(t, inst, b)
    inst, ok = GetInstance(a.Manifest.InstanceId)
    assert.Nil(t, inst)
    assert.False(t, ok)
    assert.Len(t, capabilitySuppliers[*b.Manifest.Capabilities[0]], 1)
    assert.Len(t, capabilityDependents[*b.Manifest.Dependencies[0]], 1)
}
