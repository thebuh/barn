package main

import (
	"bytes"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func LoadTestConfig() {

	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")

	// any approach to require this configuration into your program.
	var yamlExample = []byte(`
monitors:
  http:
    remote:
      name: "Some remote url"
      description: "Some remote url description"
      url: http://127.0.0.1/test
    remote2:
      name: "Some remote url2"
      description: "Some remote url2 description"
      url: http://127.0.0.2/test
  file:
    local:
      name: "Local file"
      description: "Some local file description"
      path: /tmp/test
  dummy:
    fake:
      is_safe: true
`)

	viper.ReadConfig(bytes.NewBuffer(yamlExample))
}

func TestBarnServer_AddMonitorTest(t *testing.T) {
	var barn = New()
	sm := NewSafetyMonitorDummy("dummy", "name", "description", false)
	barn.AddMonitor(sm)
	assert.Equal(t, barn.monitors[sm.GetId()], sm, "List of monitors should contain monitor")
	assert.Equal(t, "dummy", barn.GetMonitor(sm.GetId()).GetId(), "List of monitors should contain monitor")
}

func TestBarnServer_RemoveMonitorTest(t *testing.T) {
	var barn = New()
	sm := NewSafetyMonitorDummy("dummy", "name", "description", false)
	barn.AddMonitor(sm)
	barn.RemoveMonitor("dummy")
	assert.Nil(t, barn.monitors["dummy"], "should be nil")
}
func TestBarnServer_GetMonitorByIndex(t *testing.T) {
	var barn = New()
	sm := NewSafetyMonitorDummy("dummy", "name", "description", false)
	barn.AddMonitor(sm)
	sm2 := NewSafetyMonitorDummy("dummy2", "name", "description", true)
	barn.AddMonitor(sm2)
	smr, _ := barn.GetMonitorByIndex(0)
	assert.Equal(t, "dummy", smr.GetId(), "should be equal")
	assert.Equal(t, false, smr.IsSafe(), "should be equal")
	smr, _ = barn.GetMonitorByIndex(1)
	assert.Equal(t, "dummy2", smr.GetId(), "should be equal")
	_, err := barn.GetMonitorByIndex(2)
	assert.Error(t, err, "should be error")
}

func TestBarnServer_LoadConfig(t *testing.T) {
	LoadTestConfig()
	var barn = New()
	mCfg := viper.GetViper()
	barn.LoadMonitorsFromConfig(mCfg)
	assert.NotNil(t, barn.GetMonitor("remote"), "shouldn't be nil")
	assert.NotNil(t, barn.GetMonitor("remote2"), "shouldn't be nil")
	switch v := barn.GetMonitor("remote").(type) {
	case *SafetyMonitorHttp:
		assert.Equal(t, "http://127.0.0.1/test", v.url, "should be equal")
		assert.Equal(t, "Some remote url", v.GetName(), "should be equal")
		assert.Equal(t, "Some remote url description", v.GetDescription(), "should be equal")
	default:
		assert.Fail(t, "Wrong type")
	}
	assert.NotNil(t, barn.GetMonitor("local"), "shouldn't be nil")
	switch v := barn.GetMonitor("local").(type) {
	case *SafetyMonitorFile:
		assert.Equal(t, v.path, "/tmp/test", "should be equal")
	default:
		assert.Fail(t, "Wrong type")
	}
	assert.NotNil(t, barn.GetMonitor("fake"), "shouldn't be nil")
}
