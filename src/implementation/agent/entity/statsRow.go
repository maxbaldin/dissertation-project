package entity

import (
	"bytes"
	"crypto/md5"
	"fmt"
)

//easyjson:json
type StatsRow struct {
	hash    string
	Process Process
	Packet  Packet
}

func (s *StatsRow) Hash() string {
	if s.hash == "" {
		var buff bytes.Buffer

		buff.WriteString(fmt.Sprintf("%v", s.Process.Id))
		buff.WriteString(s.Process.Name)
		buff.WriteString(s.Process.Path)
		buff.WriteString(fmt.Sprintf("%v", s.Process.Sender))

		buff.WriteString(s.Packet.SourceIp)
		buff.WriteString(s.Packet.TargetIp)
		buff.WriteString(fmt.Sprintf("%v", s.Packet.TargetPort))
		buff.WriteString(fmt.Sprintf("%v", s.Packet.SourcePort))

		s.hash = fmt.Sprintf("%x", md5.Sum(buff.Bytes()))
	}

	return s.hash
}
