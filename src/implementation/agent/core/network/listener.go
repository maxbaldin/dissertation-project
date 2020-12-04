package network

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/entity"
	log "github.com/sirupsen/logrus"
)

type Listener struct {
	transformer PacketTransformer
	producer    PacketProducer
}

type PacketProducer interface {
	Produce(packet entity.StatsRow)
}

type PacketTransformer interface {
	Transform(packet gopacket.Packet) (entity.StatsRow, error)
}

func NewListener(transformer PacketTransformer, producer PacketProducer) *Listener {
	return &Listener{transformer: transformer, producer: producer}
}

func (l *Listener) Listen(handle *pcap.Handle) {
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		statsRow, err := l.transformer.Transform(packet)
		if err != nil {
			log.Info(err)
		}
		l.producer.Produce(statsRow)
	}
}
