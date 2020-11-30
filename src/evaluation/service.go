package evaluation

import (
	"bufio"
	"context"
	"math/rand"
	"net"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	ReconnectTimeout      = time.Second * 3
	MaxBytesPerConnection = 2048
)

type TestService struct {
	dependencies []OutboundDependency
	listenAddr   string
	logger       *log.Entry
}

type OutboundDependency struct {
	Addr                           string `yaml:"addr"`
	Protocol                       string `yaml:"protocol"`
	PacketSize                     int    `yaml:"packet_size"`
	DurationSeconds                int    `yaml:"duration_in_seconds"`
	TimeBetweenPacketsMilliseconds int    `yaml:"time_between_packets_ms"`
}

func NewTestService(deps []OutboundDependency, listenAddr string, logger *log.Entry) *TestService {
	return &TestService{dependencies: deps, listenAddr: listenAddr, logger: logger}
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
		ts.logger.Warnf("Unable connect to the target service %s (%s): reconnecting", dependency.Addr, err)
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
		ts.logger.Debugf("Sending data to %s:%s", dependency.Protocol, dependency.Addr)
		_, err := connBuffer.Write(packet)
		if err != nil {
			ts.logger.Warnf("Unable write to the target service %s (%s): reconnecting", dependency.Addr, err)
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
		ts.logger.Fatalf("Error listening: %s", err)
	}
	defer l.Close()
	ts.logger.Infof("Listen on %s", ts.listenAddr)

	for {
		if needToStop {
			break
		}
		c, err := l.Accept()
		if err != nil {
			ts.logger.Warnf("Unable to accept connection %s", err)
			continue
		}
		go ts.handleConnection(c)
	}
}

func (ts *TestService) handleConnection(conn net.Conn) {
	buf := make([]byte, MaxBytesPerConnection)
	_, err := conn.Read(buf)
	if err != nil {
		ts.logger.Warnf("Error reading: %s", err)
		return
	}
	err = conn.Close()
	if err != nil {
		ts.logger.Warnf("Error closing: %s", err)
	}
	ts.logger.Debug("Connection handled")
}
