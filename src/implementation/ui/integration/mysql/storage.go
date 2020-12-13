package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/maxbaldin/dissertation-project/src/implementation/ui/entity"
)

const LastHoursCnt = 2

const graphDataSQL = `SELECT outbound.hostname     as from_hostname,
       outbound.process_name as from_service,
       inbound.hostname      as to_hostname,
       inbound.process_name  as to_service,
       SUM(outbound.packets) as packetsCnt,
       SUM(outbound.size)    as sizeBytes
FROM inbound_traffic inbound
         JOIN outbound_traffic outbound
              ON inbound.target_port = outbound.target_port AND
                 inbound.target_ip = outbound.target_ip AND
                 inbound.source_port = outbound.source_port AND
                 inbound.source_ip = outbound.source_ip AND
                 inbound.hour = outbound.hour
                  AND inbound.target_ip > 0 AND outbound.target_ip > 0
                  AND cast(concat(inbound.date, ' ', inbound.hour + ':00:00') as datetime) >= now() - INTERVAL %d HOUR
                  AND cast(concat(outbound.date, ' ', outbound.hour + ':00:00') as datetime) >= now() - INTERVAL %d HOUR
GROUP BY inbound.hostname, inbound.process_name, outbound.hostname, outbound.process_name
ORDER BY inbound.process_name, outbound.process_name;`

type Storage struct {
	mysql              *sql.DB
	graphDataStatement *sql.Stmt
}

func NewNodesStorage(connectionString string) (*Storage, error) {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	graphDataStmt, err := db.Prepare(fmt.Sprintf(graphDataSQL, LastHoursCnt, LastHoursCnt))
	if err != nil {
		return nil, err
	}

	return &Storage{
		mysql:              db,
		graphDataStatement: graphDataStmt,
	}, nil
}

func (s *Storage) GraphData() (elements entity.GraphDataElements, err error) {
	rows, err := s.graphDataStatement.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var element entity.GraphDataElement
		err := rows.Scan(&element.SourceHost, &element.SourceService, &element.TargetHost, &element.TargetService, &element.PacketsCnt, &element.SizeBytes)
		if err != nil {
			return nil, err
		}
		elements = append(elements, element)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return elements, nil
}

func (s *Storage) Close() {
	_ = s.mysql.Close()
}
