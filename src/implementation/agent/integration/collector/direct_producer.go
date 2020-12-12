package collector

import (
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

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

	log.Debug("Init aggregator")
	go producer.aggregator.Aggregate(producer.inQueue, producer.aggregatedQueue)

	log.Debug("Init producer")
	go producer.produce()

	return producer
}

func (p *DirectProducer) produce() {
	for statsRow := range p.aggregatedQueue {
		b, err := statsRow.MarshalJSON()
		if err != nil {
			log.Warnf("Stats row marshaling error %s", err)
			continue
		}
		_, err = http.PostForm(p.collectorAddr, url.Values{
			"entity": {string(b)},
		})
		if err != nil {
			log.Warnf("Unable to post data to the collector %s", err)
		}
	}
}

func (p *DirectProducer) Produce(packet entity.StatsRow) {
	p.inQueue <- packet
}

func (p *DirectProducer) Close() {
	close(p.inQueue)
	close(p.aggregatedQueue)
}
