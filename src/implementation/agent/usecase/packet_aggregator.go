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

func NewAggregator(flushInterval time.Duration, startBuffLength int) *Aggregator {
	aggregator := &Aggregator{
		buffer:      make(map[string]*entity.StatsRow, startBuffLength),
		flushTicker: time.NewTicker(flushInterval),
	}

	go func() {
		for range aggregator.flushTicker.C {
			if aggregator.outChan == nil {
				continue
			}
			aggregator.mu.Lock()
			if len(aggregator.buffer) > 0 {
				log.Printf("Flushing aggregates (%d elements)", len(aggregator.buffer))
				for _, aggregatedRow := range aggregator.buffer {
					aggregator.outChan <- *aggregatedRow
				}
				log.Println("Flushing aggregates finished")
				oldBuffLen := len(aggregator.buffer)
				aggregator.buffer = make(map[string]*entity.StatsRow, oldBuffLen/2)
			} else {
				log.Println("No aggregates to flush")
			}
			aggregator.mu.Unlock()
		}
	}()

	return aggregator
}

func (a *Aggregator) Aggregate(in chan entity.StatsRow, out chan entity.StatsRow) {
	a.outChan = out
	for row := range in {
		a.mu.Lock()
		rowHash := row.Hash()
		if _, exist := a.buffer[rowHash]; exist {
			a.buffer[rowHash].Packet.Size = a.buffer[rowHash].Packet.Size + row.Packet.Size
			a.buffer[rowHash].Packet.Packets = a.buffer[rowHash].Packet.Packets + 1
		} else {
			rowPtr := row
			a.buffer[rowHash] = &rowPtr
		}
		a.mu.Unlock()
	}
}

func (a *Aggregator) Close() {
	a.flushTicker.Stop()
}
