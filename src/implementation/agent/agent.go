package main

import (
	"net"
	"sync"
	"time"

	"github.com/google/gopacket/pcap"
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/integration/collector"
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/usecase"
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/usecase/network"
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/usecase/process"
	log "github.com/sirupsen/logrus"
)

func main() {
	interfaces, err := net.Interfaces()
	checkErr(err)

	nodesRepository, err := collector.NewNodeRepository("http://collector/known_nodes", time.Second)
	checkErr(err)

	processRepository, err := process.NewRepository(nodesRepository, time.Microsecond)
	checkErr(err)
	defer processRepository.Close()

	transformer := usecase.NewPacketTransformer(processRepository)

	aggregator := usecase.NewAggregator(time.Second*5, 1000)
	producer := collector.NewDirectProducer("http://collector", aggregator, 1000)
	listener := network.NewListener(transformer, producer)

	var wg sync.WaitGroup
	for _, networkInterface := range interfaces {
		wg.Add(1)
		handle, err := pcap.OpenLive(networkInterface.Name, 65536, true, pcap.BlockForever)
		if err != nil {
			log.Warn(err)
			continue
		}
		err = handle.SetBPFFilter("tcp")
		checkErr(err)
		go listener.Listen(handle)
	}
	wg.Wait()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
