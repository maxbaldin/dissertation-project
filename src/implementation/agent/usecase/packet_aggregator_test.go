package usecase_test

import (
	"testing"
	"time"

	"github.com/maxbaldin/dissertation-project/src/implementation/agent/entity"
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/usecase"
	"github.com/stretchr/testify/assert"
)

func TestAggregator_Aggregate(t *testing.T) {
	aggTime := time.Second * 1
	aggregator := usecase.NewAggregator(aggTime, 10)
	inChan := make(chan entity.StatsRow, 10)

	row1 := entity.StatsRow{
		Process: &entity.Process{
			Id:                         12,
			Name:                       "one",
			Path:                       "/var/path",
			Sender:                     true,
			CommunicationWithKnownNode: true,
		},
		Packet: &entity.Packet{
			SourceIp:   "66.0.0.1",
			SourcePort: 1233,
			TargetIp:   "88.434.1.23",
			TargetPort: 322,
			Size:       100,
			Packets:    1,
		},
	}
	row1.Hash()

	row2 := entity.StatsRow{
		Process: &entity.Process{
			Id:                         43,
			Name:                       "two",
			Path:                       "/var/path",
			Sender:                     false,
			CommunicationWithKnownNode: false,
		},
		Packet: &entity.Packet{
			SourceIp:   "55.2.12.4",
			SourcePort: 534,
			TargetIp:   "77.3.1.55",
			TargetPort: 7878,
			Size:       1,
			Packets:    1,
		},
	}
	row2.Hash()

	inChan <- row1
	inChan <- row1
	inChan <- row1
	inChan <- row1
	inChan <- row2

	close(inChan)

	outChan := make(chan entity.StatsRow, 10)
	go aggregator.Aggregate(inChan, outChan)

	var cnt int
	for newRow := range outChan {
		if newRow.Process.Name == "one" {
			cnt++
			row1.Packet.Packets = 4
			row1.Packet.Size = 400
			assert.Equal(t, row1, newRow)
		} else if newRow.Process.Name == "two" {
			cnt++
			assert.Equal(t, row2, newRow)
		} else {
			t.Fatal("unknown stat row")
		}
		if cnt == 2 {
			break
		}
	}
}
