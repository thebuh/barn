package app

import (
	"errors"
	"fmt"
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thebuh/barn/internal/monitor"
	"github.com/thebuh/barn/internal/weather"
)

type Server interface {
	GetMonitorIds() []string
	GetMonitor(Id string) monitor.SafetyMonitor
	GetMonitorByIndex(index int) (monitor.SafetyMonitor, error)
	GetWeatherIds() []string
	GetWeather(Id string) weather.ObservingConditions
	GetWeatherByIndex(index int) (weather.ObservingConditions, error)
}

type server struct {
	monitors map[string]monitor.SafetyMonitor
	weather  map[string]weather.ObservingConditions
}

func New() *server {
	var server = server{}
	server.monitors = make(map[string]monitor.SafetyMonitor)
	server.weather = make(map[string]weather.ObservingConditions)
	return &server
}

func (s *server) LoadMonitorsFromConfig(v *viper.Viper) {
	http := v.GetStringMap("monitors.http")
	if http != nil {
		for id := range http {
			vt := v.Sub(fmt.Sprintf("monitors.http.%s", id))
			rule := monitor.NewSafetyMatchingRule(vt.GetBool("rule.invert"), vt.GetString("rule.pattern"))
			sm := monitor.NewSafetyMonitorHttp(id, vt.GetString("name"), vt.GetString("description"), vt.GetString("url"), rule)
			s.AddMonitor(sm)
		}
	}
	file := v.GetStringMap("monitors.file")
	if file != nil {
		for id := range file {
			vt := v.Sub(fmt.Sprintf("monitors.file.%s", id))
			rule := monitor.NewSafetyMatchingRule(vt.GetBool("rule.invert"), vt.GetString("rule.pattern"))
			sm := monitor.NewSafetyMonitorFile(id, vt.GetString("name"), vt.GetString("description"), vt.GetString("path"), rule)
			s.AddMonitor(sm)
		}
	}
	dummy := v.GetStringMap("monitors.dummy")
	if dummy != nil {
		for id := range dummy {
			vt := v.Sub(fmt.Sprintf("monitors.dummy.%s", id))
			sm := monitor.NewSafetyMonitorDummy(id, vt.GetString("name"), vt.GetString("description"), vt.GetBool("is_safe"))
			s.AddMonitor(sm)
		}
	}
}

func (s *server) LoadWeatherFromConfig(v *viper.Viper) {
	weatherConfig := v.GetStringMap("weather.dummy")
	if weatherConfig != nil {
		for id := range weatherConfig {
			vt := v.Sub(fmt.Sprintf("weather.dummy.%s", id))
			wt := weather.NewObservingConditionsDummy(id, vt.GetString("name"), vt.GetString("description"))
			s.AddWeather(wt)
		}
	}
	weatherConfig = v.GetStringMap("weather.http")
	if weatherConfig != nil {
		for id := range weatherConfig {
			vt := v.Sub(fmt.Sprintf("weather.http.%s", id))
			wt, _ := weather.NewObservingConditionsHttp(id, vt.GetString("name"), vt.GetString("description"), vt.GetString("url"))
			s.AddWeather(wt)
		}
	}
}

func (s *server) AddWeather(weather weather.ObservingConditions) {
	s.weather[weather.GetId()] = weather
}

func (s *server) AddMonitor(mon monitor.SafetyMonitor) {
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
	sort.Sort(sort.StringSlice(keys))
	return keys
}

func (s *server) GetMonitor(id string) monitor.SafetyMonitor {
	return s.monitors[id]
}

func (s *server) GetMonitorByIndex(id int) (monitor.SafetyMonitor, error) {
	keys := s.GetMonitorIds()
	if id > len(keys)-1 || id < 0 {
		return nil, errors.New("Index out of range")
	}
	return s.monitors[keys[id]], nil
}

func (s *server) GetWeatherIds() []string {
	keys := make([]string, 0)
	for key := range s.weather {
		keys = append(keys, key)
	}
	sort.Sort(sort.StringSlice(keys))
	return keys
}

func (s *server) GetWeather(id string) weather.ObservingConditions {
	return s.weather[id]
}

func (s *server) GetWeatherByIndex(id int) (weather.ObservingConditions, error) {
	keys := s.GetWeatherIds()
	if id > len(keys)-1 || id < 0 {
		return nil, errors.New("Index out of range")
	}
	return s.weather[keys[id]], nil
}

func (s *server) Refresh() {
	for _, val := range s.monitors {
		m := val
		go func() {
			m.Refresh()
			log.WithFields(log.Fields{
				"monitor": m.GetName(),
				"state":   m.IsSafe(),
			}).Info(fmt.Sprintf("[BARN] Monitor [%s]. Refreshing state. Now: [%t]", m.GetName(), m.IsSafe()))
		}()
	}
	for _, val := range s.weather {
		w := val
		go func() {
			w.Refresh()
			log.WithFields(log.Fields{
				"weather": w.GetName(),
				"state":   w.GetState(),
			}).Info(fmt.Sprintf("[BARN] Weather [%s]. Refreshing state.", w.GetName()))
		}()
	}
}
