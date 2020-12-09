package collector_test

import (
	"log"
	"testing"
	"time"

	"github.com/maxbaldin/dissertation-project/src/implementation/agent/entity"
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/integration/collector"
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/usecase"
)

func TestDirectProducer_Produce(t *testing.T) {
	aggregator := usecase.NewAggregator(time.Second, 1000)
	producer := collector.NewDirectProducer("http://127.0.0.1:8080", aggregator, 1000)

	for {
		log.Println("Produce!")
		producer.Produce(entity.StatsRow{
			Hostname: "test",
			Process: &entity.Process{
				Id:                         0,
				Name:                       "",
				Path:                       "",
				Sender:                     false,
				CommunicationWithKnownNode: false,
			},
			Packet: &entity.Packet{
				SourceIp:   "",
				SourcePort: 0,
				TargetIp:   "",
				TargetPort: 0,
				Size:       0,
				Packets:    0,
			},
		})
	}
}
