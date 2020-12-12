package repository

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type KnownNodesDB interface {
	KnownNodes() ([]string, error)
}

type KnownNodesRepository struct {
	mu           sync.RWMutex
	updateTicker *time.Ticker
	nodes        []string
	db           KnownNodesDB
}

func NewKnownNodesRepository(db KnownNodesDB, updateInterval time.Duration) (*KnownNodesRepository, error) {
	ticker := time.NewTicker(updateInterval)
	repo := &KnownNodesRepository{updateTicker: ticker, db: db}

	go func(t *time.Ticker) {
		for range t.C {
			err := repo.update()
			if err != nil {
				log.Warn(err)
			}
		}
	}(ticker)

	return repo, repo.update()
}

func (k *KnownNodesRepository) update() error {
	nodes, err := k.db.KnownNodes()
	if err != nil {
		return err
	}
	k.mu.Lock()
	k.nodes = nodes
	k.mu.Unlock()
	return nil
}

func (k *KnownNodesRepository) Get() []string {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.nodes
}
