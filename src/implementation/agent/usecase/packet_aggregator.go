package usecase

import (
	"log"
	"sync"
	"time"

	"github.com/maxbaldin/dissertation-project/src/implementation/agent/entity"
)

type Aggregator struct {
	minSendCnt  int
	buffer      map[string]*entity.StatsRow
	flushTicker *time.Ticker
	outChan     chan entity.StatsRow
	mu          sync.Mutex
}

func NewAggregator(flushInterval time.Duration, startBuffLength int, minSendCnt int) *Aggregator {
	aggregator := &Aggregator{
		minSendCnt:  minSendCnt,
		buffer:      make(map[string]*entity.StatsRow, startBuffLength),
		flushTicker: time.NewTicker(flushInterval),
	}

	go func() {
		for range aggregator.flushTicker.C {
			if aggregator.outChan == nil {
				continue
			}
			log.Println("Flushing aggregates")
			if len(aggregator.buffer) < aggregator.minSendCnt {
				continue
			}
			aggregator.mu.Lock()
			for _, aggregatedRow := range aggregator.buffer {
				aggregator.outChan <- *aggregatedRow
			}
			oldBuffLen := len(aggregator.buffer)
			aggregator.buffer = make(map[string]*entity.StatsRow, oldBuffLen)
			aggregator.mu.Unlock()
		}
	}()

	return aggregator
}

func (a *Aggregator) Aggregate(in chan entity.StatsRow, out chan entity.StatsRow) {
	a.outChan = out
	for row := range in {
		a.mu.Lock()
		log.Println("Aggregating..")
		rowHash := row.Hash()
		if _, exist := a.buffer[rowHash]; exist {
			a.buffer[rowHash].Packet.Size = row.Packet.Size
			a.buffer[rowHash].Packet.Packets += 1
		} else {
			a.buffer[rowHash] = &row
		}
		a.mu.Unlock()
	}
}

func (a *Aggregator) Close() {
	a.flushTicker.Stop()
}
