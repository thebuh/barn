package main

import (
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"time"
)

type SafetyMonitor interface {
	IsSafe() bool
	Refresh()
	GetId() string
	GetName() string
	GetDescription() string
}

func IsSafeString(content string) bool {
	if strings.HasPrefix(strings.ToLower(content), "true") || strings.HasPrefix(strings.ToLower(content), "1") {
		return true
	}
	return false
}

type SafetyMonitorHttp struct {
	id              string
	name            string
	description     string
	safe            bool
	url             string
	LastRefreshTime time.Time
}

func NewSafetyMonitorHttp(id string, name string, description string, url string) *SafetyMonitorHttp {
	monitor := &SafetyMonitorHttp{id: id, name: name, description: description, url: url}
	monitor.Refresh()
	return monitor
}
func NewSafetyMonitorHttpFromCfg(id string, cfg map[string]string) *SafetyMonitorHttp {
	h := &SafetyMonitorHttp{id: id, name: cfg["name"], description: cfg["description"], url: cfg["url"], safe: false}
	h.Refresh()
	return h
}

func (sm *SafetyMonitorHttp) GetId() string {
	return sm.id
}

func (sm *SafetyMonitorHttp) GetName() string {
	return sm.name
}

func (sm *SafetyMonitorHttp) GetDescription() string {
	return sm.description
}

func (sm *SafetyMonitorHttp) IsSafe() bool {
	return sm.safe
}

func (sm *SafetyMonitorHttp) Refresh() {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	response, err := client.Get(sm.url)
	if err != nil {
		sm.safe = false
		return
	}
	buf := make([]byte, 1024)
	_, err = response.Body.Read(buf)
	if err != nil && err != io.EOF {
		sm.safe = false
		return
	}
	_ = response.Body.Close()
	content := string(buf)
	sm.safe = IsSafeString(content)
	sm.LastRefreshTime = time.Now()
}

func NewSafetyMonitorDummyFromCfg(id string, cfg map[string]string) *SafetyMonitorDummy {
	dummy := &SafetyMonitorDummy{id: id, name: cfg["name"], description: cfg["description"]}
	if cfg["is_safe"] != "true" {
		dummy.safe = false
	} else {
		dummy.safe = true
	}
	return dummy
}

func NewSafetyMonitorDummy(id string, name string, description string, isSafe bool) *SafetyMonitorDummy {
	return &SafetyMonitorDummy{
		id:          id,
		name:        name,
		description: description,
		safe:        isSafe,
	}
}

type SafetyMonitorDummy struct {
	safe        bool
	id          string
	name        string
	description string
}

func (sm *SafetyMonitorDummy) IsSafe() bool {
	return sm.safe
}

func (sm *SafetyMonitorDummy) Refresh() {
}

func (sm *SafetyMonitorDummy) GetId() string {
	return sm.id
}

func (sm *SafetyMonitorDummy) GetName() string {
	return sm.name
}
func (sm *SafetyMonitorDummy) GetDescription() string {
	return sm.description
}

type SafetyMonitorFile struct {
	id              string
	name            string
	description     string
	safe            bool
	path            string
	LastRefreshTime time.Time
}

func (sm *SafetyMonitorFile) GetId() string {
	return sm.id
}

func (sm *SafetyMonitorFile) GetName() string {
	return sm.name
}
func (sm *SafetyMonitorFile) GetDescription() string {
	return sm.description
}

func (sm *SafetyMonitorFile) IsSafe() bool {
	return sm.safe
}

func (sm *SafetyMonitorFile) Refresh() {
	f, err := os.OpenFile(sm.path, os.O_RDONLY, 0444)
	if errors.Is(err, fs.ErrNotExist) {
		sm.safe = false
		return
	}
	defer f.Close()
	buf := make([]byte, 1024)
	_, err = f.Read(buf)
	if err != nil {
		sm.safe = false
		return
	}
	content := string(buf)
	sm.safe = IsSafeString(content)
	sm.LastRefreshTime = time.Now()
}

func NewSafetyMonitorFileFromCfg(id string, cfg map[string]string) *SafetyMonitorFile {
	file := &SafetyMonitorFile{id: id, name: cfg["name"], description: cfg["description"], path: cfg["path"]}
	file.Refresh()
	return file
}

func NewSafetyMonitorFile(id string, name string, description string, path string) *SafetyMonitorFile {
	file := &SafetyMonitorFile{id: id, name: name, description: description, path: path}
	file.Refresh()
	return file
}
