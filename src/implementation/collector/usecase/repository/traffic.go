package repository

import (
	"github.com/maxbaldin/dissertation-project/src/implementation/collector/entity"
	"github.com/maxbaldin/dissertation-project/src/implementation/collector/usecase/utils"
)

const YMDFormat = "2006-01-02"

type TrafficDB interface {
	InsertOrUpdateTraffic(inbound bool, date string, hour, minute int, processName, hostname string, sourceIp, sourcePort, targetIp, targetPort int, packets, size uint) error
}

type TrafficRepository struct {
	db TrafficDB
}

func NewTrafficRepository(db TrafficDB) *TrafficRepository {
	return &TrafficRepository{
		db: db,
	}
}

func (tr *TrafficRepository) Persist(row entity.Traffic) error {
	return tr.db.InsertOrUpdateTraffic(
		row.Inbound,
		row.Date.Format(YMDFormat),
		row.Date.Hour(),
		row.Date.Minute(),
		row.ProcessName,
		row.Hostname,
		int(utils.Ip2int(row.SourceIP)),
		row.SourcePort,
		int(utils.Ip2int(row.TargetIp)),
		row.TargetPort,
		row.PacketsCnt,
		row.Size,
	)
}
