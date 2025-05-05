package main

import (
	"barn/cmd"
	"bytes"
	"github.com/spf13/viper"
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
  dummy:
    fake:
      is_safe: true
`)

	viper.ReadConfig(bytes.NewBuffer(yamlExample))
}

func main() {
	viper.SetConfigName("barn")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()
	//fakeConfig()

	//if err != nil {
	//	panic(fmt.Errorf("fatal error config file: %w", err))
	//}
	barnSrv := New()
	mCfg := viper.GetViper()
	barnSrv.LoadMonitorsFromConfig(mCfg)
	discovery := NewDiscoverySever(32227, 8080)
	api := NewApiServer(barnSrv, 8080)
	go discovery.Start()
	defer discovery.Close()
	go api.Start()
	for {
		barnSrv.Refresh()
		time.Sleep(30)
	}
	cmd.Execute()
}
