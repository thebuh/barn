package monitor

import (
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"regexp"
	"time"
)

type SafetyMonitor interface {
	IsSafe() bool
	Refresh()
	GetId() string
	GetName() string
	GetDescription() string
	GetRawValue() string
	GetTimeStamp() time.Time
}

type SafetyMatchingRule struct {
	invert  bool
	pattern string

	regex *regexp.Regexp
}

func NewSafetyMatchingRule(invert bool, pattern string) *SafetyMatchingRule {
	rule := &SafetyMatchingRule{
		invert:  false,
		pattern: pattern,
	}
	regex, err := regexp.Compile("(?i)" + pattern)
	if err != nil || pattern == "" {
		regex, _ = regexp.Compile(`(?i)true|1`)
	}
	rule.regex = regex
	return rule
}

func (rule *SafetyMatchingRule) isSafe(content string) bool {
	if rule.regex == nil {
		return false
	}
	if rule.invert {
		return !rule.regex.MatchString(content)
	}
	return rule.regex.MatchString(content)
}

type SafetyMonitorHttp struct {
	id              string
	name            string
	description     string
	safe            bool
	url             string
	lastRefreshTime time.Time
	lastValue       string
	rule            *SafetyMatchingRule
	client          *http.Client
}

func NewSafetyMonitorHttp(id string, name string, description string, url string, rule *SafetyMatchingRule) *SafetyMonitorHttp {
	monitor := &SafetyMonitorHttp{id: id, name: name, description: description, url: url}
	monitor.client = &http.Client{
		Timeout: 5 * time.Second,
	}
	monitor.rule = rule
	monitor.Refresh()
	return monitor
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

func (sm *SafetyMonitorHttp) GetRawValue() string {
	return sm.lastValue
}

func (sm *SafetyMonitorHttp) GetTimeStamp() time.Time {
	return sm.lastRefreshTime
}

func (sm *SafetyMonitorHttp) Refresh() {
	response, err := sm.client.Get(sm.url)
	if err != nil {
		sm.safe = false
		sm.lastValue = ""
		return
	}
	buf := make([]byte, 1024)
	n, err := response.Body.Read(buf)
	if err != nil && err != io.EOF {
		sm.safe = false
		sm.lastValue = ""
		return
	}
	buf = buf[:n]
	_ = response.Body.Close()
	content := string(buf)
	sm.lastValue = content
	sm.safe = sm.rule.isSafe(content)
	sm.lastRefreshTime = time.Now()
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
func (sm *SafetyMonitorDummy) GetRawValue() string {
	return ""
}

func (sm *SafetyMonitorDummy) GetTimeStamp() time.Time {
	return time.Time{}
}

type SafetyMonitorFile struct {
	id              string
	name            string
	description     string
	safe            bool
	path            string
	lastRefreshTime time.Time
	lastValue       string
	rule            *SafetyMatchingRule
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

func (sm *SafetyMonitorFile) GetRawValue() string {
	return sm.lastValue
}

func (sm *SafetyMonitorFile) GetTimeStamp() time.Time {
	return sm.lastRefreshTime
}

func (sm *SafetyMonitorFile) IsSafe() bool {
	return sm.safe
}

func (sm *SafetyMonitorFile) Refresh() {
	f, err := os.OpenFile(sm.path, os.O_RDONLY, 0444)
	if errors.Is(err, fs.ErrNotExist) {
		sm.safe = false
		sm.lastValue = ""
		return
	}
	defer f.Close()
	buf := make([]byte, 1024)
	n, err := f.Read(buf)
	if err != nil {
		sm.safe = false
		sm.lastValue = ""
		return
	}
	buf = buf[:n]
	content := string(buf)
	sm.lastValue = content
	sm.safe = sm.rule.isSafe(content)
	sm.lastRefreshTime = time.Now()
}

func NewSafetyMonitorFile(id string, name string, description string, path string, rule *SafetyMatchingRule) *SafetyMonitorFile {
	file := &SafetyMonitorFile{id: id, name: name, description: description, path: path, rule: rule}
	file.Refresh()
	return file
}
