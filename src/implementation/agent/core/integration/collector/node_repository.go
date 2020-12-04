package collector

import (
	"log"
	"time"
)

type NodesRepository struct {
	updateTicker *time.Ticker
}

func NewNodeRepository(nodesListAddr string, updateInterval time.Duration) (*NodesRepository, error) {
	ticker := time.NewTicker(updateInterval)
	repo := &NodesRepository{updateTicker: ticker}

	go func(t *time.Ticker) {
		for range t.C {
			err := repo.update()
			if err != nil {
				log.Println(err)
			}
		}
	}(ticker)

	return repo, repo.update()
}

func (nr *NodesRepository) update() error {
	return nil
}

func (nr *NodesRepository) IsKnownNode(ip string) bool {
	return false
}

func (nr *NodesRepository) Close() {
	nr.updateTicker.Stop()
}
