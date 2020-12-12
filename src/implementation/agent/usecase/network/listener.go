package network

import (
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/entity"
	log "github.com/sirupsen/logrus"
)

const StateReportingInterval = time.Second * 5

type Listener struct {
	transformErrors int
	transformOKs    int
	transformer     PacketTransformer
	producer        PacketProducer
}

type PacketProducer interface {
	Produce(packet entity.StatsRow)
}

type PacketTransformer interface {
	Transform(packet gopacket.Packet) (entity.StatsRow, error)
}

func NewListener(transformer PacketTransformer, producer PacketProducer) *Listener {
	listener := &Listener{transformer: transformer, producer: producer}
	go func() {
		for range time.NewTicker(StateReportingInterval).C {
			log.Infof(
				"Transformation errors %d, Success transformations %d",
				listener.transformErrors,
				listener.transformOKs,
			)
		}
	}()
	return listener
}

func (l *Listener) Listen(handle *pcap.Handle, interfaceName string, maxLen int) {
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		statsRow, err := l.transformer.Transform(packet)
		if err != nil {
			l.transformErrors += 1
			continue
		}
		l.transformOKs += 1
		l.producer.Produce(statsRow)
	}
	log.Infof("End of listening of the network interface %s", interfaceName)
}
