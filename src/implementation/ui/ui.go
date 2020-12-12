package main

import (
	"net/http"
	"os"

	"github.com/maxbaldin/dissertation-project/src/implementation/ui/controller"
	"github.com/maxbaldin/dissertation-project/src/implementation/ui/entity"
	"github.com/maxbaldin/dissertation-project/src/implementation/ui/integration/mysql"
	"github.com/maxbaldin/dissertation-project/src/implementation/ui/usecase"
	"github.com/maxbaldin/dissertation-project/src/implementation/ui/usecase/repository"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Initializing...")

	cfg := config()

	http.Handle("/", http.FileServer(http.Dir(cfg.Server.ServeFolder)))

	nodesDb, err := mysql.NewNodesStorage(cfg.Integration.Db.ConnectionString)
	if err != nil {
		panic(err)
	}
	nodeTransformer := usecase.NewNodeTransformer()
	nodeRepository := repository.NewNodesRepository(nodesDb, nodeTransformer)
	apiController := controller.NewApi(nodeRepository)

	http.HandleFunc("/api/graph", apiController.Handle)

	log.Info("Application is ready")
	log.Fatal(http.ListenAndServe(cfg.Server.ListenAddr, nil))

}

func config() *entity.Config {
	var cfg entity.Config
	if cfgPath, exist := os.LookupEnv("UI_CFG"); !exist {
		panic("you pass path to the config file in 'UI_CFG' environment variable")
	} else {
		err := cfg.FromFile(cfgPath)
		if err != nil {
			panic(err)
		}
		return &cfg
	}
}
