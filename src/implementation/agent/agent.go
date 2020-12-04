package main

import (
	"net"
	"sync"
	"time"

	"github.com/google/gopacket/pcap"
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/core"
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/core/integration/collector"
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/core/network"
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/core/process"
)

func main() {
	interfaces, err := net.Interfaces()
	checkErr(err)

	nodesRepository, err := collector.NewNodeRepository("", time.Second*5)
	checkErr(err)

	processRepository, err := process.NewRepository(nodesRepository, time.Second*5)
	checkErr(err)
	defer processRepository.Close()

	transformer := core.NewPacketTransformer(processRepository)
	aggregator := core.NewAggregator(time.Second*1, 100)
	producer := collector.NewDirectProducer(aggregator, 1000)
	listener := network.NewListener(transformer, producer)

	var wg sync.WaitGroup
	for _, networkInterface := range interfaces {
		wg.Add(1)
		handle, err := pcap.OpenLive(networkInterface.Name, 65536, true, pcap.BlockForever)
		checkErr(err)
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
