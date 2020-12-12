package controller

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/maxbaldin/dissertation-project/src/implementation/collector/entity"
)

const (
	StateReportingInterval = time.Second * 5
	SuccessResponseReply   = "ok"
)

type TrafficRepository interface {
	Persist(row entity.Traffic) error
}

type Transformer interface {
	Transform(r *http.Request) (entity.Traffic, error)
}

type CollectController struct {
	collectedItems int
	errorItems     int
	repository     TrafficRepository
	transformer    Transformer
}

func NewCollectController(repository TrafficRepository, transformer Transformer) *CollectController {
	controller := &CollectController{
		repository:  repository,
		transformer: transformer,
	}

	go func() {
		for range time.NewTicker(StateReportingInterval).C {
			log.Infof(
				"Collected %d, Rejected %d",
				controller.collectedItems,
				controller.errorItems,
			)
		}
	}()

	return controller
}

func (c *CollectController) Handle(w http.ResponseWriter, r *http.Request) {
	trafficEntity, err := c.transformer.Transform(r)
	if err != nil {
		c.errorItems += 1
		log.Warn(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = c.repository.Persist(trafficEntity)
	if err != nil {
		c.errorItems += 1
		log.Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.collectedItems += 1
	_, _ = fmt.Fprint(w, SuccessResponseReply)
}
