package evaluation

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

const ReconnectTimeout = time.Second * 3

type TestService struct {
	dependencies []OutboundDependency
	listenAddr   string
}

type OutboundDependency struct {
	Addr                           string `yaml:"addr"`
	Protocol                       string `yaml:"protocol"`
	PacketSize                     int    `yaml:"packet_size"`
	DurationSeconds                int    `yaml:"duration_in_seconds"`
	TimeBetweenPacketsMilliseconds int    `yaml:"time_between_packets_ms"`
}

func NewTestService(deps []OutboundDependency, listenAddr string) *TestService {
	return &TestService{dependencies: deps, listenAddr: listenAddr}
}

func (ts *TestService) Run(ctx context.Context) {
	var wg sync.WaitGroup
	if ts.listenAddr != "" {
		wg.Add(1)
		go ts.listen(ctx, &wg)
	}
	for _, dep := range ts.dependencies {
		wg.Add(1)

		concreteCtx := ctx
		if dep.DurationSeconds > 0 {
			concreteCtx, _ = context.WithTimeout(ctx, time.Second*time.Duration(dep.DurationSeconds))
		}
		go ts.writer(dep, concreteCtx, &wg)
	}
	wg.Wait()
}

func (ts *TestService) writer(dependency OutboundDependency, ctx context.Context, group *sync.WaitGroup) {
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
		log.Println("Unable connect to the target service " + dependency.Addr + ": reconnecting")
		time.Sleep(ReconnectTimeout)
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

func (ts *TestService) listen(ctx context.Context, wg *sync.WaitGroup) {
	var needToStop bool
	defer wg.Done()

	go func() {
		<-ctx.Done()
		needToStop = true
	}()

	l, err := net.Listen("tcp4", ts.listenAddr)
	if err != nil {
		log.Fatal("Error listening:", err.Error())
	}
	defer l.Close()
	log.Println("Listen on", ts.listenAddr)

	for {
		if needToStop {
			break
		}
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Unable to accept connection", err)
			return
		}
		go ts.handleConnection(c)
	}
}

func (ts *TestService) handleConnection(conn net.Conn) {
	buf := make([]byte, 2048)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	_, err = conn.Write([]byte("Message received."))
	if err != nil {
		fmt.Println("Error writing:", err.Error())
	}
	err = conn.Close()
	if err != nil {
		fmt.Println("Error closing:", err.Error())
	}
	log.Println("Connection handled")
}
