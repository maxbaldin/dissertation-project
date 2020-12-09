package usecase

import (
	"errors"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/entity"
)

var (
	ErrTCPLayerIsNil       = errors.New("tcp layer is nil")
	ErrIPLayerIsNil        = errors.New("ip layer is nil")
	ErrUnableToFindProcess = errors.New("unable to find process")

	hostname string
)

type ProcessRepository interface {
	FindByNetworkActivity(packet entity.Packet) (proc entity.Process)
}

type PacketTransformer struct {
	processRepository ProcessRepository
}

func init() {
	hostname, _ = os.Hostname()
}

func NewPacketTransformer(procRepository ProcessRepository) *PacketTransformer {
	return &PacketTransformer{processRepository: procRepository}
}

func (pt *PacketTransformer) Transform(packet gopacket.Packet) (statsRow entity.StatsRow, err error) {
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		sourcePort := int(tcp.SrcPort)
		targetPort := int(tcp.DstPort)

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
			if process.Id < 0 { // invalid process
				return statsRow, ErrUnableToFindProcess
			}

			// inbound: source(remote):target(local)
			// outbound: source(local):target(remote)
			// remove extra information in case of unknown remotes
			if !process.CommunicationWithKnownNode {
				if process.Sender {
					packet.TargetIp = "unk"
					packet.TargetPort = 0
					packet.SourcePort = 0
				} else {
					packet.SourceIp = "unk"
					packet.SourcePort = 0
					packet.TargetPort = 0
				}
			}

			return entity.StatsRow{
				Process:  &process,
				Packet:   &packet,
				Hostname: hostname,
			}, nil
		} else {
			return statsRow, ErrIPLayerIsNil
		}
	} else {
		return statsRow, ErrTCPLayerIsNil
	}
}
