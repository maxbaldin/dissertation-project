package usecase

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"time"

	agentEntity "github.com/maxbaldin/dissertation-project/src/implementation/agent/entity"
	"github.com/maxbaldin/dissertation-project/src/implementation/collector/entity"
)

var (
	ErrEmptyEntityField = errors.New("empty entity field")
)

type RequestTransformer struct {
}

func NewRequestTransformer() *RequestTransformer {
	return &RequestTransformer{}
}

func (rt *RequestTransformer) Transform(r *http.Request) (e entity.Traffic, err error) {
	input := r.PostFormValue("entity")
	if input == "" {
		return e, ErrEmptyEntityField
	}
	var row agentEntity.StatsRow
	err = json.Unmarshal([]byte(input), &row)
	if err != nil {
		return e, err
	}

	e.Inbound = row.Process.Sender == false
	e.Date = time.Now()
	e.ProcessName = row.Process.Name
	e.Hostname = row.Hostname

	e.SourceIP = net.ParseIP(row.Packet.SourceIp)
	if e.SourceIP == nil {
		e.SourceIP = net.ParseIP("0.0.0.0")
	}
	e.SourcePort = row.Packet.SourcePort

	e.TargetIp = net.ParseIP(row.Packet.TargetIp)
	if e.TargetIp == nil {
		e.TargetIp = net.ParseIP("0.0.0.0")
	}
	e.TargetPort = row.Packet.TargetPort

	e.PacketsCnt = row.Packet.Packets
	e.Size = row.Packet.Size

	return e, nil
}
