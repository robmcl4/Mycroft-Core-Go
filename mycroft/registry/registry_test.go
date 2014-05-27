package registry


import (
    "testing"
    "github.com/robmcl4/mycroft/app"
    "github.com/coreos/go-semver/semver"
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
    cap1.Version, _ = semver.NewVersion("1.1.1")
    cap1.Name = "capability1"
    caps[0] = cap1
    cap2 := new(app.Capability)
    cap2.Version, _ = semver.NewVersion("1.1.2")
    cap2.Name = "capability2"
    caps[1] = cap2
    a.Manifest.Capabilities = caps

    // construct dependencies
    deps := make([]*app.Capability, 2)
    dep1 := new(app.Capability)
    dep1.Version, _ = semver.NewVersion("0.1.2")
    dep1.Name = "dependency1"
    deps[0] = dep1
    dep2 := new(app.Capability)
    dep2.Version, _ = semver.NewVersion("0.1.3")
    dep2.Name = "dependency2"
    deps[1] = dep2

    a.Manifest.Dependencies = deps

    return a
}


func setUp() {
    capabilitySuppliers = make(map[string]map[semver.Version][]*app.App)
    capabilityDependents = make(map[string]map[semver.Version][]*app.App)
    instances = make(map[string]*app.App)
}


func TestCapabilitySuppliersMap(t *testing.T) {
    setUp()
    a := getMockApp()
    Register(a)
    // make sure the registry entries exist
    if _, ok := capabilitySuppliers["capability1"]; !ok {
        t.Error("Capability Suppliers did not contain 'capability1' name")
    }
    expectedVer, _ := semver.NewVersion("1.1.1")
    val, ok := capabilitySuppliers["capability1"][*expectedVer]
    if !ok {
        t.Error("Capability Suppliers did not contain '1.1.1' semver")
    }
    if len(val) != 1 {
        t.Error("Too many references to app")
    }
    if val[0] != a {
        t.Error("Capability Suppliers contains incorrect reference to app")
    }
    expectedVer, _ = semver.NewVersion("1.1.2")
    val, ok = capabilitySuppliers["capability2"][*expectedVer]
    if !ok {
        t.Error("Capability Suppliers did not contain '1.1.2' semver")
    }
    if len(val) != 1 {
        t.Error("Too many references to app (second capability)")
    }
    if val[0] != a {
        t.Error("Capability Suppliers contains incorrect reference to app")
    }
}


func TestCapabilityDependentsMap(t *testing.T) {
    setUp()
    a := getMockApp()
    Register(a)
    // make sure the registry entries exist
    if _, ok := capabilityDependents["dependency1"]; !ok {
        t.Error("Capability Suppliers did not contain 'dependency1' name")
    }
    expectedVer, _ := semver.NewVersion("0.1.2")
    val, ok := capabilityDependents["dependency1"][*expectedVer]
    if !ok {
        t.Error("Capability Suppliers did not contain '0.1.2' semver")
    }
    if len(val) != 1 {
        t.Error("Too many references to app")
    }
    if val[0] != a {
        t.Error("Capability Dependents contains incorrect reference to app")
    }
    if _, ok := capabilityDependents["dependency2"]; !ok {
        t.Error("Capability Suppliers did not contain 'dependency2' name")
    }
    expectedVer, _ = semver.NewVersion("0.1.3")
    val, ok = capabilityDependents["dependency2"][*expectedVer]
    if !ok {
        t.Error("Capability Dependents did not contain '0.1.3' semver")
    }
    if len(val) != 1 {
        t.Error("Too many references to app (second dependency)")
    }
    if val[0] != a {
        t.Error("Capability Dependents contains incorrect reference to app")
    }
}


func TestInstanceIdMap(t *testing.T) {
    setUp()
    a := getMockApp()
    Register(a)
    if retrieved, ok := instances["instance1"]; !ok {
        t.Error("App was not found in instances")
    } else if retrieved != a {
        t.Error("Did not retrieve correct app")
    } 
}


func TestGetProviders(t *testing.T) {
    setUp()
    a := getMockApp()
    Register(a)
    apps := GetProviders(a.Manifest.Capabilities[0])
    if len(apps) != 1 {
        t.Fatal("Incorrect length of apps")
    }
    if apps[0] != a {
        t.Error("retrieved incorrect app instance")
    }
    // make sure it actually returns a _copy_ of the array
    apps[0] = nil
    apps = GetProviders(a.Manifest.Capabilities[0])
    if apps[0] != a {
        t.Error("retrieved incorrect app instance after setting copy to nil")
    }
}


func TestGetDependents(t *testing.T) {
    setUp()
    a := getMockApp()
    Register(a)
    apps := GetDependents(a.Manifest.Dependencies[0])
    if len(apps) != 1 {
        t.Fatal("Incorrect length of apps")
    }
    if apps[0] != a {
        t.Error("retrieved incorrect app instance")
    }
    // make sure it actually returns a _copy_ of the array
    apps[0] = nil
    apps = GetDependents(a.Manifest.Dependencies[0])
    if apps[0] != a {
        t.Error("retrieved incorrect app instance after setting copy to nil")
    }
}


func TestGetInstance(t *testing.T) {
    setUp()
    a := getMockApp()
    Register(a)
    inst, ok := GetInstance("instance1")
    if !ok {
        t.Error("should have found app instance")
    }
    if inst != a {
        t.Error("found incorrect instance")
    }
    inst, ok = GetInstance("foo")
    if ok {
        t.Error("should not have found an instance")
    }
    if inst != nil {
        t.Error("should not have found a non-nil instance")
    }
}


func TestRemoveApp(t *testing.T) {
    setUp()
    a := getMockApp()
    b := getMockApp()
    b.Manifest.InstanceId = "foo11223"
    Register(a)
    Register(b)
    Remove(a)
    if _, ok := GetInstance(b.Manifest.InstanceId); !ok {
        t.Error("Should contain instance of b")
    }
    if _, ok := GetInstance(a.Manifest.InstanceId); ok {
        t.Error("Should not contain instance of b")
    }
    ver, _ := semver.NewVersion("1.1.1")
    apps := capabilitySuppliers["capability1"][*ver]
    if len(apps) != 1 {
        t.Error("apps is incorrect length")
    }
    if apps[0] != b {
        t.Error("returned incorrect app")
    }
}
