package evaluation

import (
	"bufio"
	"context"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

type TestService struct {
	dependencies []OutboundDependency
}

type OutboundDependency struct {
	Addr                           string `yaml:"addr"`
	Protocol                       string `yaml:"protocol"`
	PacketSize                     int    `yaml:"packet_size"`
	DurationSeconds                int    `yaml:"duration_in_seconds"`
	TimeBetweenPacketsMilliseconds int    `yaml:"time_between_packets_ms"`
}

func NewTestService(deps []OutboundDependency) *TestService {
	return &TestService{dependencies: deps}
}

func (ts *TestService) Run(ctx context.Context) {
	var wg sync.WaitGroup
	for _, dep := range ts.dependencies {
		wg.Add(1)

		concreteCtx := ctx
		if dep.DurationSeconds > 0 {
			concreteCtx, _ = context.WithTimeout(ctx, time.Second*time.Duration(dep.DurationSeconds))
		}
		go ts.worker(dep, concreteCtx, &wg)
	}
	wg.Wait()
}

func (ts *TestService) worker(dependency OutboundDependency, ctx context.Context, group *sync.WaitGroup) {
	defer func() {
		group.Done()
	}()

	var needToStop bool

	go func() {
		<-ctx.Done()
		needToStop = true
	}()

start:
	if needToStop {
		return
	}
	conn, err := net.Dial(dependency.Protocol, dependency.Addr)
	if err != nil {
		log.Println("Unable connect to the target service: reconnecting")
		goto start
	}
	connBuffer := bufio.NewWriter(conn)

	packet := make([]byte, dependency.PacketSize)
	for {
		if needToStop {
			return
		}
		rand.Read(packet)
		log.Printf("Sending data to %s:%s", dependency.Protocol, dependency.Addr)
		_, err := connBuffer.Write(packet)
		if err != nil {
			log.Println("Unable write to the target service: reconnecting")
			goto start
		}
		time.Sleep(time.Millisecond * time.Duration(dependency.TimeBetweenPacketsMilliseconds))
	}
}
