package main

import (
	"net"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/google/gopacket/pcap"
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/entity"
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/integration/collector"
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/usecase"
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/usecase/network"
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/usecase/process"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Initializing...")
	cfg := config()

	runtime.GOMAXPROCS(cfg.MaxCores)

	http.DefaultClient.Timeout = time.Duration(cfg.Integration.HTTP.DefaultTimeoutMs) * time.Millisecond

	log.Info("Initializing known node repository...")
	nodesRepository, err := collector.NewNodeRepository(
		cfg.Integration.Collector.KnownNodes.URL,
		time.Duration(cfg.Integration.Collector.KnownNodes.UpdateIntervalSec)*time.Second,
		cfg.Integration.Collector.KnownNodes.Additional,
	)
	checkErr(err)

	log.Info("Initializing process repository...")
	processRepository, err := process.NewRepository(
		nodesRepository,
		time.Duration(cfg.Integration.Process.Repository.UpdateIntervalMs)*time.Millisecond,
	)
	checkErr(err)
	defer processRepository.Close()

	log.Info("Initializing aggregator...")
	aggregator := usecase.NewAggregator(
		time.Duration(cfg.Integration.Collector.Aggregator.FlushIntervalSec)*time.Second,
		cfg.Integration.Collector.Aggregator.InitialBufferLength,
	)
	defer aggregator.Close()

	log.Info("Initializing packet transformer...")
	transformer := usecase.NewPacketTransformer(processRepository)

	log.Info("Initializing producer...")
	producer := collector.NewDirectProducer(
		cfg.Integration.Collector.Producer.URL,
		aggregator,
		cfg.Integration.Collector.Producer.QueueLength,
	)
	defer producer.Close()

	log.Info("Initializing network listener...")
	listener := network.NewListener(transformer, producer)

	var wg sync.WaitGroup

	interfaces, err := net.Interfaces()
	checkErr(err)

	log.Info("Start listening network interfaces...")
	for _, networkInterface := range interfaces {
		wg.Add(1)
		go func(interfaceName string, snapshotLength int) {
			log.Infof("Attaching to %s", interfaceName)
			handle, err := pcap.OpenLive(
				interfaceName,
				int32(snapshotLength),
				false,
				pcap.BlockForever,
			)
			if err != nil {
				log.Warn(err)
				return
			}

			err = handle.SetBPFFilter(cfg.Network.PCap.BFFFilter)
			checkErr(err)
			listener.Listen(handle, interfaceName, snapshotLength)
		}(networkInterface.Name, cfg.Network.PCap.SnapshotLength)
	}

	wg.Wait()
}

func config() *entity.Config {
	var cfg entity.Config
	if cfgPath, exist := os.LookupEnv("AGENT_CFG"); !exist {
		panic("you pass path to the config file in 'AGENT_CFG' environment variable")
	} else {
		err := cfg.FromFile(cfgPath)
		if err != nil {
			panic(err)
		}
		return &cfg
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
