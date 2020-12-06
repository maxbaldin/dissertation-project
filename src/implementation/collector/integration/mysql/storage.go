package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	mysql *sql.DB
}

func New(connectionString string) (*Storage, error) {
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		return nil, err
	}
	return &Storage{
		mysql: db,
	}, nil
}

func (s *Storage) InsertOrUpdateTraffic(inbound bool, date string, hour int, processName, hostname string, sourceIp, sourcePort, targetIp, targetPort int, packets, size uint) error {
	return nil
}

func (s *Storage) KnownNodes(sinceDate string) ([]string, error) {
	return []string{}, nil
}

func (s *Storage) Close() {
	_ = s.mysql.Close()
}
