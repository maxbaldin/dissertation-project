package entity

type Packet struct {
	Process Process

	SourceIp   string
	SourcePort int

	TargetIp   string
	TargetPort int

	Size int
}
