package main

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/maxbaldin/dissertation-project/src/evaluation"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Service              Service
	OutboundDependencies []evaluation.OutboundDependency `yaml:"outbound_dependencies"`
}

type Service struct {
	ListenAddr string `yaml:"listen_addr"`
	Name       string
}

func main() {
	var config Config
	if cfgPath, exist := os.LookupEnv("CFG"); !exist {
		panic("you pass path to the config file in 'cfgPath' environment variable")
	} else {
		cfgData, err := ioutil.ReadFile(cfgPath)
		checkErr(err)
		checkErr(yaml.Unmarshal(cfgData, &config))
	}

	service := evaluation.NewTestService(config.OutboundDependencies, config.Service.ListenAddr)

	service.Run(context.Background())
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
