package process

import (
	"log"
	"sync"
	"time"

	"github.com/maxbaldin/dissertation-project/src/implementation/agent/entity"
)

var UnknownProcess = entity.Process{
	Id:   -1,
	Name: "Unknown",
	Path: "Unknown",
}

type Repository struct {
	updateTicker     *time.Ticker
	mu               sync.RWMutex
	connectionsIndex NetTCPRowsIndex
	processesIndex   ProcessesIndex
}

func NewRepository(updateInterval time.Duration) (*Repository, error) {
	ticker := time.NewTicker(updateInterval)
	repo := &Repository{updateTicker: ticker}

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

func (r *Repository) update() error {
	connections, err := NewNetTCPIndex()
	if err != nil {
		return err
	}
	processes, err := GetProcMap()
	if err != nil {
		return err
	}

	r.mu.Lock()
	r.connectionsIndex = connections
	r.processesIndex = processes
	r.mu.Unlock()

	return nil
}

func (r *Repository) FindByNetworkActivity(packet entity.Packet) (proc entity.Process) {
	srcInfo, srcInfoExist := r.connectionsIndex.LookupSource(packet.SourceIp, packet.SourcePort)
	targetInfo, targetInfoExist := r.connectionsIndex.LookupSource(packet.TargetIp, packet.TargetPort)

	if srcInfoExist {
		if process, exist := r.processesIndex[srcInfo.Inode]; exist {
			proc.Id = process.PID
			proc.Name = process.Name
			proc.Sender = true
			return proc
		}
	} else if targetInfoExist {
		if process, exist := r.processesIndex[targetInfo.Inode]; exist {
			proc.Id = process.PID
			proc.Name = process.Name
			return proc
		}
	}

	return UnknownProcess
}

func (r *Repository) Close() {
	r.updateTicker.Stop()
}
