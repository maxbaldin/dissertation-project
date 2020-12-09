package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const inboundTrafficTable = "inbound_traffic"
const outboundTrafficTable = "outbound_traffic"

const insertOrUpdateStatement = `INSERT INTO collector.%s (date, 
                                       hour,
                                       process_name,
                                       hostname,
                                       source_ip,
                                       source_port,
                                       target_ip,
                                       target_port,
                                       packets,
                                       size)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE packets = packets + VALUES(packets),
                        size    = size + VALUES(size);`

const knownNodesStatement = `SELECT GROUP_CONCAT(ip) as ips FROM known_nodes`

type Storage struct {
	mysql                      *sql.DB
	inboundInsertOrUpdateStmt  *sql.Stmt
	outboundInsertOrUpdateStmt *sql.Stmt
	knownNodesStmt             *sql.Stmt
}

func New(connectionString string) (*Storage, error) {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	inboundInsertOrUpdateStmt, err := db.Prepare(fmt.Sprintf(insertOrUpdateStatement, inboundTrafficTable))
	if err != nil {
		return nil, err
	}

	outboundInsertOrUpdateStmt, err := db.Prepare(fmt.Sprintf(insertOrUpdateStatement, outboundTrafficTable))
	if err != nil {
		return nil, err
	}

	knownNodesStmt, err := db.Prepare(knownNodesStatement)
	if err != nil {
		return nil, err
	}

	return &Storage{
		mysql:                      db,
		inboundInsertOrUpdateStmt:  inboundInsertOrUpdateStmt,
		outboundInsertOrUpdateStmt: outboundInsertOrUpdateStmt,
		knownNodesStmt:             knownNodesStmt,
	}, nil
}

func (s *Storage) InsertOrUpdateTraffic(inbound bool, date string, hour int, processName, hostname string, sourceIp, sourcePort, targetIp, targetPort int, packets, size uint) error {
	var currStmt *sql.Stmt
	if inbound {
		currStmt = s.inboundInsertOrUpdateStmt
	} else {
		currStmt = s.outboundInsertOrUpdateStmt
	}

	_, err := currStmt.Exec(date, hour, processName, hostname, sourceIp, sourcePort, targetIp, targetPort, packets, size)
	return err
}

func (s *Storage) KnownNodes() ([]string, error) {
	row := s.knownNodesStmt.QueryRow()
	var ids sql.NullString
	err := row.Scan(&ids)
	if err != nil {
		return nil, err
	}

	return strings.Split(ids.String, ","), nil
}

func (s *Storage) Close() {
	_ = s.mysql.Close()
}
