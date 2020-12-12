package main

import (
	"net/http"
	"os"
	"time"

	"github.com/maxbaldin/dissertation-project/src/implementation/collector/entity"
	log "github.com/sirupsen/logrus"

	"github.com/maxbaldin/dissertation-project/src/implementation/collector/controller"
	"github.com/maxbaldin/dissertation-project/src/implementation/collector/integration/mysql"
	"github.com/maxbaldin/dissertation-project/src/implementation/collector/usecase"
	"github.com/maxbaldin/dissertation-project/src/implementation/collector/usecase/repository"
)

func main() {
	log.Info("Initializing...")

	cfg := config()

	idxController := controller.NewIndex()
	http.HandleFunc("/", idxController.Handle)

	mysqlDB, err := mysql.New(cfg.Integration.Db.ConnectionString)
	if err != nil {
		panic(err)
	}
	trafficRepo := repository.NewTrafficRepository(mysqlDB)
	reqTransformer := usecase.NewRequestTransformer()
	collectController := controller.NewCollectController(trafficRepo, reqTransformer)

	http.HandleFunc("/collect", collectController.Handle)

	knownNodesRepository, err := repository.NewKnownNodesRepository(
		mysqlDB,
		time.Duration(cfg.KnownNodes.UpdateIntervalSec)*time.Second,
	)
	if err != nil {
		panic(err)
	}
	knownNodesController := controller.NewKnownNodesController(knownNodesRepository)
	http.HandleFunc("/known_nodes", knownNodesController.Handle)

	log.Info("Application is ready")
	log.Fatal(http.ListenAndServe(cfg.Server.ListenAddr, nil))
}

func config() *entity.Config {
	var cfg entity.Config
	if cfgPath, exist := os.LookupEnv("COLLECTOR_CFG"); !exist {
		panic("you pass path to the config file in 'COLLECTOR_CFG' environment variable")
	} else {
		err := cfg.FromFile(cfgPath)
		if err != nil {
			panic(err)
		}
		return &cfg
	}
}
