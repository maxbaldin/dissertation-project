package collector

import (
	"log"

	"github.com/maxbaldin/dissertation-project/src/implementation/agent/entity"
)

type DirectProducer struct {
	aggregator      Aggregator
	inQueue         chan entity.StatsRow
	aggregatedQueue chan entity.StatsRow
}

type Aggregator interface {
	Aggregate(in chan entity.StatsRow, out chan entity.StatsRow)
}

func NewDirectProducer(aggregator Aggregator, queueLen int) *DirectProducer {
	producer := &DirectProducer{
		aggregator:      aggregator,
		inQueue:         make(chan entity.StatsRow, queueLen),
		aggregatedQueue: make(chan entity.StatsRow, queueLen),
	}

	go producer.handle()

	go producer.produce()

	return producer
}

func (p *DirectProducer) produce() {
	for statsRow := range p.aggregatedQueue {
		_, err := statsRow.MarshalJSON()
		if err != nil {
			log.Println(err)
		}
		log.Printf("Need to send %v", statsRow)
	}
}

func (p *DirectProducer) handle() {
	p.aggregator.Aggregate(p.inQueue, p.aggregatedQueue)
}

func (p *DirectProducer) Produce(packet entity.StatsRow) {
	p.inQueue <- packet
}

func (p *DirectProducer) Close() {
	close(p.inQueue)
	close(p.aggregatedQueue)
}
