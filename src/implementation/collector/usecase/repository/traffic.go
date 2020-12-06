package repository

import (
	"log"

	"github.com/maxbaldin/dissertation-project/src/implementation/collector/entity"
	"github.com/maxbaldin/dissertation-project/src/implementation/collector/usecase/utils"
)

const YMDFormat = "02-01-2006"

type TrafficDB interface {
	InsertOrUpdateTraffic(inbound bool, date string, hour int, processName, hostname string, sourceIp, sourcePort, targetIp, targetPort int, packets, size uint) error
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
	log.Println(row)
	return tr.db.InsertOrUpdateTraffic(
		row.Inbound,
		row.Date.Format(YMDFormat),
		row.Date.Hour(),
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
