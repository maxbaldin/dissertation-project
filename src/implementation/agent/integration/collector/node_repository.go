package collector

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
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
			log.Info("Updating known nodes...")
			err := repo.update()
			if err != nil {
				log.Warn(err)
			}
			log.Info("Know nodes updated...")
		}
	}(repo.updateTicker)

	return repo, repo.update()
}

func (nr *NodesRepository) update() error {
	resp, err := http.Get(nr.nodesListAddr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	nr.nodesList = strings.Split(string(b), ",")
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
