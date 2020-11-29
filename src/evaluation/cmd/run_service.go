package main

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/maxbaldin/dissertation-project/src/evaluation"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	OutboundDependencies []evaluation.OutboundDependency `yaml:"outbound_dependencies"`
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
	service := evaluation.NewTestService(config.OutboundDependencies)

	service.Run(context.Background())
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
