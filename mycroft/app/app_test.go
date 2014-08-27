package app

import (
    "github.com/stretchr/testify/assert"
    "github.com/coreos/go-semver/semver"
    "testing"
    "net"
    "sync"
    "time"
)

type mockConnection struct {
    toRead, written []byte
}
func (c *mockConnection) Read(b []byte) (int, error) {
    return copy(c.toRead, b), nil
}
func (c *mockConnection) Write(b []byte) (int, error) {
    c.written = append(c.written, b...)
    return len(b), nil
}
func (c *mockConnection) Close() (error) {return nil}
func (c *mockConnection) LocalAddr() (net.Addr) {return nil}
func (c *mockConnection) RemoteAddr() (net.Addr) {return nil}
func (c *mockConnection) SetDeadline(time.Time) (error) {return nil}
func (c *mockConnection) SetReadDeadline(time.Time) (error) {return nil}
func (c *mockConnection) SetWriteDeadline(time.Time) (error) {return nil}


func TestHasStatusConnected(t *testing.T) {
    assert.True(t, STATUS_CONNECTED >= 0, "should have STATUS_CONNECTED > 0")

    app := new(App)
    app.Status = STATUS_CONNECTED
    assert.Equal(t, app.StatusString(), "connected")
}


func TestHasStatusUp(t *testing.T) {
    assert.True(t, STATUS_UP >= 0, "should have STATUS_UP > 0")

    app := new(App)
    app.Status = STATUS_UP
    assert.Equal(t, app.StatusString(), "up")
}


func TestHasStatusDown(t *testing.T) {
    assert.True(t, STATUS_DOWN >= 0, "should have STATUS_DOWN > 0")

    app := new(App)
    app.Status = STATUS_DOWN
    assert.Equal(t, app.StatusString(), "down")
}


func TestHasStatusInUse(t *testing.T) {
    assert.True(t, STATUS_IN_USE >= 0, "should have STATUS_IN_USE > 0")

    app := new(App)
    app.Status = STATUS_IN_USE
    app.Priority = 10
    assert.Equal(t, app.StatusString(), "in_use 10")
}


func TestCapabilityFields(t *testing.T) {
    cap := new(Capability)

    v, _ := semver.NewVersion("1.1.2")
    cap.Version = *v

    cap.Name = "foo"
}


func TestManifestFields(t *testing.T) {
    m := new(Manifest)

    m.Name = "foo_name"
    m.DisplayName = "foo display name"
    m.InstanceId = "c5996e7c-518e-4683-a3b7-16357025d78f"
    m.ApiVersion = 1
    m.Description = "a totally cool test app"

    v, _ := semver.NewVersion("1.1.2")
    m.Version = *v
    m.Capabilities = make([]*Capability, 0)
    m.Dependencies = make([]*Capability, 0)
}


func TestAppFields(t *testing.T) {
    a := new(App)

    a.Connection = &net.IPConn{}
    a.Manifest = &Manifest{}
    a.Status = 1
    a.Priority = 10
    a.RWMutex = sync.RWMutex{}
}


func TestNewApp(t *testing.T) {
    a := NewApp()
    assert.Equal(t, a.Status, STATUS_CONNECTED)

    a.RWMutex.Lock()
    a.RWMutex.Unlock()
}


func TestSendWithoutManifest(t *testing.T) {
    a := NewApp()
    toSend := make(map[string]interface{})
    fakeConn := &mockConnection{}
    a.Connection = fakeConn
    a.Send("FOO", toSend)
    assert.Equal(
        t,
        fakeConn.written,
        []byte("6\nFOO {}"))
}


func TestSendWithManifest(t *testing.T) {
    a := NewApp()
    a.Manifest = &Manifest{}
    toSend := make(map[string]interface{})
    fakeConn := &mockConnection{}
    a.Connection = fakeConn
    a.Send("FOO", toSend)
    assert.Equal(
        t,
        fakeConn.written,
        []byte("6\nFOO {}"))
}


func TestSendWithBody(t *testing.T) {
    a := NewApp()
    toSend := make(map[string]interface{})
    toSend["foo"] = "bar"
    fakeConn := &mockConnection{}
    a.Connection = fakeConn
    a.Send("FOO", toSend)
    assert.Equal(
        t,
        fakeConn.written,
        []byte("17\nFOO {\"foo\":\"bar\"}"))
}
