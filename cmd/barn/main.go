package main

import (
	"bytes"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	api "github.com/thebuh/barn/internal/api"
	"github.com/thebuh/barn/internal/app"
	"github.com/thebuh/barn/pkg/discovery"
)

func fakeConfig() {
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")

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
weather:
  dummy:
    name: "Weather"
    description: "Weather conditions"
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
	barnApp := app.New()
	mCfg := viper.GetViper()
	barnApp.LoadMonitorsFromConfig(mCfg)
	barnApp.LoadWeatherFromConfig(mCfg)

	apiPort := viper.GetUint32("api.port")
	discoveryPort := viper.GetUint32("discovery.port")
	disc := discovery.NewDiscoverySever(discoveryPort, apiPort)
	api := api.NewApiServer(barnApp, apiPort)
	go disc.Start()
	defer disc.Close()
	go api.Start()
	for {
		barnApp.Refresh()
		time.Sleep(10 * time.Second)
	}
}
