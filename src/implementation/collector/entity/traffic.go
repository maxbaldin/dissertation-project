package entity

import (
	"net"
	"time"
)

type Traffic struct {
	Inbound     bool
	Date        time.Time
	ProcessName string
	Hostname    string
	SourceIP    net.IP
	SourcePort  int
	TargetIp    net.IP
	TargetPort  int
	PacketsCnt  uint
	Size        uint
}
