package core

import (
	"log"
	"time"

	"github.com/maxbaldin/dissertation-project/src/implementation/agent/entity"
)

type Aggregator struct {
	maxBuffLen  int
	buffer      map[string]*entity.StatsRow
	flushTicker *time.Ticker
}

func NewAggregator(flushInterval time.Duration, startBuffLength int) *Aggregator {
	return &Aggregator{
		buffer:      make(map[string]*entity.StatsRow, startBuffLength),
		flushTicker: time.NewTicker(flushInterval),
	}
}

func (a *Aggregator) Aggregate(in chan entity.StatsRow, out chan entity.StatsRow) {
	for row := range in {
		select {
		case <-a.flushTicker.C:
			log.Println("Flushing aggregates")
			for _, aggregatedRow := range a.buffer {
				out <- *aggregatedRow
			}
			oldBuffLen := len(a.buffer)
			a.buffer = make(map[string]*entity.StatsRow, oldBuffLen)
		default:
			rowHash := row.Hash()
			if _, exist := a.buffer[rowHash]; exist {
				a.buffer[rowHash].Packet.Size = row.Packet.Size
				a.buffer[rowHash].Packet.Packets += 1
			} else {
				a.buffer[rowHash] = &row
			}
		}
	}
}
