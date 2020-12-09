package collector

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type NodesRepository struct {
	nodesListAddr string
	nodesList     []string
	updateTicker  *time.Ticker
}

func NewNodeRepository(nodesListAddr string, updateInterval time.Duration) (*NodesRepository, error) {
	repo := &NodesRepository{updateTicker: time.NewTicker(updateInterval), nodesListAddr: nodesListAddr}

	go func(t *time.Ticker) {
		for range t.C {
			log.Println("Updating known nodes...")
			err := repo.update()
			if err != nil {
				log.Println(err)
			}
		}
	}(repo.updateTicker)

	return repo, repo.update()
}

func (nr *NodesRepository) update() error {
	resp, err := http.Get(nr.nodesListAddr)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	nr.nodesList = strings.Split(string(b), ",")
	log.Printf("Found %d known nodes", len(nr.nodesList))
	return nil
}

func (nr *NodesRepository) IsKnownNode(ip string) bool {
	for _, node := range nr.nodesList {
		if node == ip {
			return true
		}
	}
	return false
}

func (nr *NodesRepository) Close() {
	nr.updateTicker.Stop()
}
