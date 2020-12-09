package main

import (
	"net"
	"net/http"
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
	log.Println("Initializing...")

	http.DefaultClient.Timeout = time.Millisecond * 10

	nodesRepository, err := collector.NewNodeRepository("http://collector/known_nodes", time.Second)
	checkErr(err)

	processRepository, err := process.NewRepository(nodesRepository, time.Microsecond)
	checkErr(err)
	defer processRepository.Close()

	transformer := usecase.NewPacketTransformer(processRepository)

	aggregator := usecase.NewAggregator(time.Second*1, 100000)
	producer := collector.NewDirectProducer("http://collector", aggregator, 100000)
	listener := network.NewListener(transformer, producer)

	var wg sync.WaitGroup

	interfaces, err := net.Interfaces()
	checkErr(err)

	for _, networkInterface := range interfaces {
		wg.Add(1)
		go func(interfaceName string, snapshotLength int) {
			log.Printf("Attaching to %s", interfaceName)
			handle, err := pcap.OpenLive(interfaceName, int32(snapshotLength), false, pcap.BlockForever)
			if err != nil {
				log.Warn(err)
				return
			}

			err = handle.SetBPFFilter("tcp")
			checkErr(err)
			listener.Listen(handle, interfaceName, snapshotLength)
		}(networkInterface.Name, 1024*1024*128)
	}
	wg.Wait()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
