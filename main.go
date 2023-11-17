package main

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thebuh/barn/cmd"
	"time"
)

func fakeConfig() {
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")

	// any approach to require this configuration into your program.
	var yamlExample = []byte(`
monitors:
  file:
    local:
      name: "Local file"
      description: "Some local file description"
      path: /tmp/test
      rule:
		regex: '^[A-Z]+\\.com$
		invert: true
  dummy:
    fake:
      name: "Fake monitor"
      is_safe: true
`)

	viper.ReadConfig(bytes.NewBuffer(yamlExample))
}

func main() {
	viper.SetConfigName("barn")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetDefault("discovery.port", 32227)
	viper.SetDefault("api.port", 8080)
	viper.ReadInConfig()
	log.SetFormatter(&log.TextFormatter{})
	//fakeConfig()

	//if err != nil {
	//	panic(fmt.Errorf("fatal error config file: %w", err))
	//}
	barnSrv := New()
	mCfg := viper.GetViper()
	barnSrv.LoadMonitorsFromConfig(mCfg)

	apiPort := viper.GetUint32("api.port")
	discoveryPort := viper.GetUint32("discovery.port")
	discovery := NewDiscoverySever(discoveryPort, apiPort)
	api := NewApiServer(barnSrv, apiPort)
	go discovery.Start()
	defer discovery.Close()
	go api.Start()
	for {
		barnSrv.Refresh()
		time.Sleep(10 * time.Second)
	}
	cmd.Execute()
}
