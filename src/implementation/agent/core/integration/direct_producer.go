package integration

import (
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/entity"
)

type DirectProducer struct {
}

func NewDirectProducer() *DirectProducer {
	return &DirectProducer{}
}

func (p *DirectProducer) Produce(packet entity.StatsRow) error {
	return nil
}
