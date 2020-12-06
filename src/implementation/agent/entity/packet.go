package entity

type Packet struct {
	SourceIp   string
	SourcePort int

	TargetIp   string
	TargetPort int

	Size    uint
	Packets uint
}
