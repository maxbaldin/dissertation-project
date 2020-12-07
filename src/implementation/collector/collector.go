package main

import (
	"log"
	"net/http"
	"time"

	"github.com/maxbaldin/dissertation-project/src/implementation/collector/controller"
	"github.com/maxbaldin/dissertation-project/src/implementation/collector/integration/mysql"
	"github.com/maxbaldin/dissertation-project/src/implementation/collector/usecase"
	"github.com/maxbaldin/dissertation-project/src/implementation/collector/usecase/repository"
)

func main() {

	idxController := controller.NewIndex()
	http.HandleFunc("/", idxController.Handle)

	mysqlDB, err := mysql.New("collector:!VB3{&uC6uwA9M#P@tcp(mysql:3306)/collector")
	if err != nil {
		panic(err)
	}
	trafficRepo := repository.NewTrafficRepository(mysqlDB)
	reqTransformer := usecase.NewRequestTransformer()
	collectController := controller.NewCollectController(trafficRepo, reqTransformer)

	http.HandleFunc("/collect", collectController.Handle)

	knownNodesRepository, err := repository.NewKnownNodesRepository(mysqlDB, time.Second*5)
	if err != nil {
		panic(err)
	}
	knownNodesController := controller.NewKnownNodesController(knownNodesRepository)
	http.HandleFunc("/known_nodes", knownNodesController.Handle)

	log.Println("Application is ready")
	log.Fatal(http.ListenAndServe(":80", nil))
}
