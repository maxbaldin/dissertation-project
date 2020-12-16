package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/maxbaldin/dissertation-project/src/implementation/ui/entity"
)

const LastHoursCnt = 2

const graphDataSQL = `SELECT c.from_hostname,
       c.from_service,
       c.to_hostname,
       c.to_service,
       SUM(c.packetsCnt)    as packetsCnt,
       SUM(c.sizeBytes)     as sizeBytes,
       SUM(c.sizePerMinute) as sizePerMinute
FROM (
		 # Remote connections:
         SELECT outbound.hostname                                                             as from_hostname,
                outbound.process_name                                                         as from_service,
                inbound.hostname                                                              as to_hostname,
                inbound.process_name                                                          as to_service,
                SUM(outbound.packets)                                                         as packetsCnt,
                SUM(outbound.size)                                                            as sizeBytes,
                CAST(SUM(outbound.size) / COUNT(DISTINCT inbound.minute) as unsigned integer) as sizePerMinute
         FROM inbound_traffic inbound
                  JOIN outbound_traffic outbound
                       ON inbound.target_port = outbound.target_port AND
                          inbound.target_ip = outbound.target_ip AND
                          inbound.source_port = outbound.source_port AND
                          inbound.source_ip = outbound.source_ip AND
                          inbound.hour = outbound.hour AND
                          inbound.minute = outbound.minute AND
                          cast(concat(inbound.date, ' ', inbound.hour + ':00:00') as datetime) >=
                          now() - INTERVAL %d HOUR AND
                          cast(concat(outbound.date, ' ', outbound.hour + ':00:00') as datetime) >=
                          now() - INTERVAL %d HOUR
         WHERE inbound.target_ip > 0
           AND outbound.target_ip > 0
         GROUP BY inbound.hostname, inbound.process_name, outbound.hostname, outbound.process_name
         UNION
		 # Local-only connections:
         SELECT o.hostname                                                       as from_hostname,
                o.process_name                                                   as from_service,
                o2.hostname                                                      as to_hostname,
                o2.process_name                                                  as to_service,
                SUM(o.packets)                                                   as packetsCnt,
                SUM(o.size)                                                      as sizeBytes,
                CAST(SUM(o.size) / COUNT(DISTINCT o.minute) as unsigned integer) as sizePerMinute
         FROM outbound_traffic o
                  JOIN outbound_traffic o2 ON
                 o.source_ip = o2.source_ip AND
                 o.target_ip = o2.target_ip AND
                 o.source_port = o2.target_port AND
                 o.target_port = o2.source_port AND
                 o.hour = o2.hour AND
                 o.minute = o2.minute AND
                 cast(concat(o.date, ' ', o.hour + ':00:00') as datetime) >= now() - INTERVAL %d HOUR AND
                 cast(concat(o2.date, ' ', o2.hour + ':00:00') as datetime) >= now() - INTERVAL %d HOUR
         WHERE o.source_ip > 0
           AND o2.target_ip > 0
         GROUP BY o.hostname, o.process_name, o2.hostname, o2.process_name
         ORDER BY from_service, to_service
     ) as c
GROUP BY c.from_hostname, c.from_service, c.to_hostname, c.to_service;`

type Storage struct {
	mysql              *sql.DB
	graphDataStatement *sql.Stmt
}

func NewNodesStorage(connectionString string) (*Storage, error) {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	graphDataStmt, err := db.Prepare(fmt.Sprintf(graphDataSQL, LastHoursCnt, LastHoursCnt, LastHoursCnt, LastHoursCnt))
	if err != nil {
		return nil, err
	}

	return &Storage{
		mysql:              db,
		graphDataStatement: graphDataStmt,
	}, nil
}

func (s *Storage) GraphData() (entity.GraphDataElements, error) {
	rows, err := s.graphDataStatement.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var elements entity.GraphDataElements

	for rows.Next() {
		var element entity.GraphDataElement
		err := rows.Scan(
			&element.SourceHost,
			&element.SourceService,
			&element.TargetHost,
			&element.TargetService,
			&element.PacketsCnt,
			&element.SizeBytes,
			&element.SizeBytesPerMinute,
		)
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
