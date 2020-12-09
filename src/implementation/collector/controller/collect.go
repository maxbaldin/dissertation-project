package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/maxbaldin/dissertation-project/src/implementation/collector/entity"
)

const SuccessResponseReply = "ok"

type TrafficRepository interface {
	Persist(row entity.Traffic) error
}

type Transformer interface {
	Transform(r *http.Request) (entity.Traffic, error)
}

type CollectController struct {
	repository  TrafficRepository
	transformer Transformer
}

func NewCollectController(repository TrafficRepository, transformer Transformer) *CollectController {
	return &CollectController{
		repository:  repository,
		transformer: transformer,
	}
}

func (c *CollectController) Handle(w http.ResponseWriter, r *http.Request) {
	trafficEntity, err := c.transformer.Transform(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = c.repository.Persist(trafficEntity)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = fmt.Fprint(w, SuccessResponseReply)
}
