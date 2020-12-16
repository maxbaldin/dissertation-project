package main

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	"github.com/maxbaldin/dissertation-project/src/evaluation"
	"github.com/maxbaldin/dissertation-project/src/evaluation/usecase"

	yaml "gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
)

var logger *log.Logger

type Config struct {
	Service              Service
	OutboundDependencies []evaluation.OutboundDependency `yaml:"outbound_dependencies"`
}

type Service struct {
	ListenAddr string `yaml:"listen_addr"`
	Name       string
}

func main() {
	logger = log.New()
	logger.SetLevel(log.DebugLevel)

	var config Config
	if cfgPath, exist := os.LookupEnv("TEST_SERVICE_CFG"); !exist {
		panic("you pass path to the config file in 'cfgPath' environment variable")
	} else {
		cfgData, err := ioutil.ReadFile(cfgPath)
		checkErr(err)
		checkErr(yaml.Unmarshal(cfgData, &config))
	}

	entry := logger.WithField("app", config.Service.Name)

	listenAddr := usecase.ReplaceLocalhostWithOutboundIP(config.Service.ListenAddr)

	logger.Infof("Listening on %s", listenAddr)
	service := evaluation.NewTestService(config.OutboundDependencies, listenAddr, time.Second*120, entry)

	service.Run(context.Background())
}

func checkErr(err error) {
	if err != nil {
		logger.Fatal(err)
	}
}
