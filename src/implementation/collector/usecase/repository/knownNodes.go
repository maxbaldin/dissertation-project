package repository

import (
	"log"
	"sync"
	"time"
)

type KnownNodesDB interface {
	KnownNodes(date string) ([]string, error)
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
				log.Println(err)
			}
		}
	}(ticker)

	return repo, repo.update()
}

func (k *KnownNodesRepository) update() error {
	now := time.Now()
	nodes, err := k.db.KnownNodes(now.Format(YMDFormat))
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
