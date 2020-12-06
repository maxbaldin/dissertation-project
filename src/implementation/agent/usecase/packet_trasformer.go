package usecase

import (
	"errors"
	"strconv"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/entity"
)

var (
	ErrTCPLayerIsNil = errors.New("tcp layer is nil")
	ErrIPLayerIsNil  = errors.New("ip layer is nil")
)

type ProcessRepository interface {
	FindByNetworkActivity(packet entity.Packet) (proc entity.Process)
}

type PacketTransformer struct {
	processRepository ProcessRepository
}

func NewPacketTransformer(procRepository ProcessRepository) *PacketTransformer {
	return &PacketTransformer{processRepository: procRepository}
}

func (pt *PacketTransformer) Transform(packet gopacket.Packet) (statsRow entity.StatsRow, err error) {
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		sourcePort, _ := strconv.Atoi(tcp.SrcPort.String())
		targetPort, _ := strconv.Atoi(tcp.DstPort.String())

		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer != nil {
			ip, _ := ipLayer.(*layers.IPv4)
			sourceIp := ip.SrcIP
			targetIp := ip.DstIP
			packet := entity.Packet{
				SourceIp:   sourceIp.String(),
				SourcePort: sourcePort,
				TargetIp:   targetIp.String(),
				TargetPort: targetPort,
				Size:       uint(packet.Metadata().Length),
				Packets:    1,
			}
			process := pt.processRepository.FindByNetworkActivity(packet)

			// remove information about the unknown nodes
			if process.CommunicationWithKnownNode {
				if process.Sender {
					packet.TargetIp = ""
					packet.TargetPort = 0
				} else {
					packet.TargetIp = ""
					packet.SourcePort = 0
				}
			}

			return entity.StatsRow{
				Process: process,
				Packet:  packet,
			}, nil
		} else {
			return statsRow, ErrIPLayerIsNil
		}
	} else {
		return statsRow, ErrTCPLayerIsNil
	}
}
