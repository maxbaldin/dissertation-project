package collector

import (
	"log"
	"net/http"
	"net/url"

	"github.com/maxbaldin/dissertation-project/src/implementation/agent/entity"
)

type DirectProducer struct {
	collectorAddr   string
	aggregator      Aggregator
	inQueue         chan entity.StatsRow
	aggregatedQueue chan entity.StatsRow
}

type Aggregator interface {
	Aggregate(in chan entity.StatsRow, out chan entity.StatsRow)
}

func NewDirectProducer(collectorAddr string, aggregator Aggregator, queueLen int) *DirectProducer {
	producer := &DirectProducer{
		collectorAddr:   collectorAddr,
		aggregator:      aggregator,
		inQueue:         make(chan entity.StatsRow, queueLen),
		aggregatedQueue: make(chan entity.StatsRow, queueLen),
	}

	log.Println("Init handler")
	go producer.handle()

	log.Println("Init producer")
	go producer.produce()

	return producer
}

func (p *DirectProducer) produce() {
	for statsRow := range p.aggregatedQueue {
		b, err := statsRow.MarshalJSON()
		if err != nil {
			log.Println(err)
			continue
		}
		_, err = http.PostForm(p.collectorAddr+"/collect", url.Values{
			"entity": {string(b)},
		})
		if err != nil {
			log.Println(err)
		}
		log.Println("Produce", string(b))
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
