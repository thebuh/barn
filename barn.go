package main

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"sort"
)

type Server interface {
	GetMonitorIds() []string
	GetMonitor(Id string) SafetyMonitor
	GetMonitorByIndex(index int) (SafetyMonitor, error)
}

type server struct {
	monitors map[string]SafetyMonitor
}

func New() *server {
	var server = server{}
	server.monitors = make(map[string]SafetyMonitor)
	return &server
}

func (s *server) LoadMonitorsFromConfig(v *viper.Viper) {
	http := v.GetStringMap("monitors.http")
	if http != nil {
		for id := range http {
			vt := v.Sub(fmt.Sprintf("monitors.http.%s", id))
			rule := NewSafetyMatchingRule(vt.GetBool("rule.invert"), vt.GetString("rule.pattern"))
			sm := NewSafetyMonitorHttp(id, vt.GetString("name"), vt.GetString("description"), vt.GetString("url"), rule)
			s.AddMonitor(sm)
		}
	}
	file := v.GetStringMap("monitors.file")
	if file != nil {
		for id := range file {
			vt := v.Sub(fmt.Sprintf("monitors.file.%s", id))
			rule := NewSafetyMatchingRule(vt.GetBool("rule.invert"), vt.GetString("rule.pattern"))
			sm := NewSafetyMonitorFile(id, vt.GetString("name"), vt.GetString("description"), vt.GetString("path"), rule)
			s.AddMonitor(sm)
		}
	}
	dummy := v.GetStringMap("monitors.dummy")
	if dummy != nil {
		for id := range dummy {
			vt := v.Sub(fmt.Sprintf("monitors.dummy.%s", id))
			sm := NewSafetyMonitorDummy(id, vt.GetString("name"), vt.GetString("description"), vt.GetBool("is_safe"))
			s.AddMonitor(sm)
		}
	}
}

func (s *server) AddMonitor(mon SafetyMonitor) {
	s.monitors[mon.GetId()] = mon
}

func (s *server) RemoveMonitor(id string) {
	delete(s.monitors, id)
}

func (s *server) GetMonitorIds() []string {
	keys := make([]string, 0)
	for key := range s.monitors {
		keys = append(keys, key)
	}
	return keys
}

func (s *server) GetMonitor(id string) SafetyMonitor {
	return s.monitors[id]
}

func (s *server) GetMonitorByIndex(id int) (SafetyMonitor, error) {
	keySlice := make([]string, 0)
	for key := range s.monitors {
		keySlice = append(keySlice, key)
	}

	sort.Strings(keySlice)
	if id > len(keySlice)-1 || id < 0 {
		return nil, errors.New("Index out of range")
	}
	return s.monitors[keySlice[id]], nil
}

func (s *server) Refresh() {
	for _, val := range s.monitors {
		go val.Refresh()
	}
}
