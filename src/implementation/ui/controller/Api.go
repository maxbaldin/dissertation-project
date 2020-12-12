package controller

import (
	"net/http"

	"github.com/maxbaldin/dissertation-project/src/implementation/ui/entity"
	log "github.com/sirupsen/logrus"
)

type ApiResponse struct {
	Nodes entity.Nodes
	Edges entity.Edges
}

type NodesRepository interface {
	FindAll() (resp entity.GraphResponse, err error)
}

type Api struct {
	nodesRepository NodesRepository
}

func NewApi(nodesRepository NodesRepository) *Api {
	return &Api{nodesRepository: nodesRepository}
}

func (a *Api) Handle(w http.ResponseWriter, r *http.Request) {
	resp, err := a.nodesRepository.FindAll()
	if err != nil {
		log.Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b, err := resp.MarshalJSON()
	if err != nil {
		log.Warn(err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	w.Write(b)
}
