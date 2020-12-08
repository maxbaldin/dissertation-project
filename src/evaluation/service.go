package evaluation

import (
	"context"
	"io"
	"math/rand"
	"net"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	ReconnectInterval = time.Millisecond * 10
)

type TestService struct {
	dependencies      []OutboundDependency
	listenAddr        string
	listenConnTimeout time.Duration
	logger            *log.Entry
}

type OutboundDependency struct {
	Addr                           string `yaml:"addr"`
	Protocol                       string `yaml:"protocol"`
	PacketSize                     int    `yaml:"packet_size"`
	TimeBetweenPacketsMilliseconds int    `yaml:"time_between_packets_ms"`
}

func NewTestService(deps []OutboundDependency, listenAddr string, listenConnTimeout time.Duration, logger *log.Entry) *TestService {
	return &TestService{dependencies: deps, listenConnTimeout: listenConnTimeout, listenAddr: listenAddr, logger: logger}
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
	ts.logger.Info("Creating TCP connection")
	conn, err := net.Dial(dependency.Protocol, dependency.Addr)
	if err != nil {
		ts.logger.Warnf("Unable connect to the target service %s (%s): sleep %s and reconnecting", dependency.Addr, err, ReconnectInterval)
		time.Sleep(ReconnectInterval)
		goto start
	}

	packet := make([]byte, dependency.PacketSize)
	for {
		if needToStop {
			ts.logger.Infof("Need to stop writing")
			return
		}
		rand.Read(packet)
		_, err := conn.Write(packet)
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
			ts.logger.Infof("Accepting of the new connections is stopped")
			break
		}
		ts.logger.Infof("Accepting new connection")
		c, err := l.Accept()
		if err != nil {
			ts.logger.Warnf("Unable to accept connection %s", err)
			continue
		}
		if err != nil {
			ts.logger.Warnf("Unable to set connection deadline %s", err)
		}
		go ts.handleConnection(c)
	}
}

func (ts *TestService) handleConnection(conn net.Conn) {
	var bytesCnt int
	tmp := make([]byte, 256)
	startTime := time.Now()
	for {
		if time.Since(startTime) > ts.listenConnTimeout {
			ts.logger.Info("Stopping listening by timeout")
			break
		}
		n, err := conn.Read(tmp)
		bytesCnt += n
		if err != nil {
			if err != io.EOF {
				ts.logger.Warnf("Read error:", err)
			}
			break
		}
	}

	_ = conn.Close()
	ts.logger.Infof("Connection closed after %d bytes", bytesCnt)
}
