package main

import (
	"net"
	"sync"
	"time"

	"github.com/google/gopacket/pcap"
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/core"
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/core/integration"
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/core/network"
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/core/process"
)

func main() {
	interfaces, err := net.Interfaces()
	checkErr(err)

	processRepository, err := process.NewRepository(time.Second * 5)
	checkErr(err)
	defer processRepository.Close()

	transformer := core.NewPacketTransformer(processRepository)
	producer := integration.NewDirectProducer()
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
