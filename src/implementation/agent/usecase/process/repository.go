package process

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/maxbaldin/dissertation-project/src/implementation/agent/entity"
)

var UnknownProcess = entity.Process{
	Id:   -1,
	Name: "Unknown",
	Path: "Unknown",
}

type NodesRepository interface {
	IsKnownNode(ip string) bool
}

type Repository struct {
	updateTicker     *time.Ticker
	nodesRepository  NodesRepository
	connectionsIndex NetTCPRowsIndex
	processesIndex   ProcessesIndex
	mu               sync.RWMutex
}

func NewRepository(nodesRepository NodesRepository, updateInterval time.Duration) (*Repository, error) {
	ticker := time.NewTicker(updateInterval)
	repo := &Repository{updateTicker: ticker, nodesRepository: nodesRepository}

	go func(r *Repository) {
		for range r.updateTicker.C {
			connections, err := NewNetTCPIndex()
			if err != nil {
				log.Warn(err)
				continue
			}
			r.mu.Lock()
			r.connectionsIndex = connections
			r.mu.Unlock()
		}
	}(repo)

	go func(r *Repository) {
		for range r.updateTicker.C {
			processesMap, err := GetProcMap()
			if err != nil {
				log.Warn(err)
				continue
			}
			r.mu.Lock()
			r.processesIndex = processesMap
			r.mu.Unlock()
		}
	}(repo)

	return repo, repo.updateAll()
}

func (r *Repository) updateAll() error {
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
	r.mu.RLock()
	defer r.mu.RUnlock()

	srcInfo, srcInfoExist := r.connectionsIndex.LookupSource(packet.SourceIp, packet.SourcePort)
	srcProcessInfo, srcInodeExist := r.processesIndex.InodeExist(srcInfo.Inode)
	srcFit := srcInfoExist &&
		packet.SourceIp == srcInfo.Local.IpV4.String() &&
		packet.SourcePort == srcInfo.Local.Port &&
		packet.TargetIp == srcInfo.Remote.IpV4.String() &&
		packet.TargetPort == srcInfo.Remote.Port &&
		srcInodeExist

	targetInfo, targetInfoExist := r.connectionsIndex.LookupSource(packet.TargetIp, packet.TargetPort)
	targetProcessInfo, targetInodeExist := r.processesIndex.InodeExist(targetInfo.Inode)
	targetFit := targetInfoExist &&
		packet.SourceIp == targetInfo.Remote.IpV4.String() &&
		packet.SourcePort == targetInfo.Remote.Port &&
		packet.TargetIp == targetInfo.Local.IpV4.String() &&
		packet.TargetPort == targetInfo.Local.Port &&
		targetInodeExist

	if srcFit {
		proc.Id = srcProcessInfo.PID
		proc.Name = srcProcessInfo.Name
		proc.Sender = true
		proc.CommunicationWithKnownNode = r.nodesRepository.IsKnownNode(packet.TargetIp)
		return proc
	} else if targetFit {
		proc.Id = targetProcessInfo.PID
		proc.Name = targetProcessInfo.Name
		proc.CommunicationWithKnownNode = r.nodesRepository.IsKnownNode(packet.SourceIp)
		return proc
	}

	return UnknownProcess
}

func (r *Repository) Close() {
	r.updateTicker.Stop()
}
