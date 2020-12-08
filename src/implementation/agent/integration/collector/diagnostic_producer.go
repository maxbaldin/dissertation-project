package collector

import (
	"os"
	"sort"
	"strconv"
	"time"

	units "github.com/docker/go-units"
	"github.com/maxbaldin/dissertation-project/src/implementation/agent/entity"
	"github.com/olekukonko/tablewriter"
	"github.com/paulbellamy/ratecounter"
)

type Stats struct {
	ProcessName string
	Sent        *ratecounter.RateCounter
	Received    *ratecounter.RateCounter
	Packets     *ratecounter.RateCounter
}

type DiagnosticProducer struct {
	data chan entity.StatsRow
}

func NewDiagnosticProducer(chanSize int) *DiagnosticProducer {
	producer := &DiagnosticProducer{
		data: make(chan entity.StatsRow, chanSize),
	}
	go producer.Visualise()
	return producer
}

func (p *DiagnosticProducer) Produce(packet entity.StatsRow) {
	p.data <- packet
}

func (p *DiagnosticProducer) Visualise() {
	packetsStats := map[int]*Stats{}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"PID", "Process name", "Packets", "Sent", "Received"})

	go func() {
		ticker := time.Tick(time.Second * 1)
		for range ticker {
			table.ClearRows()
			keys := make([]int, 0, len(packetsStats))
			for k := range packetsStats {
				keys = append(keys, k)
			}
			sort.Ints(keys)

			for _, k := range keys {
				table.Append([]string{
					strconv.Itoa(k),
					packetsStats[k].ProcessName,
					strconv.Itoa(int(packetsStats[k].Packets.Rate())),
					units.BytesSize(float64(packetsStats[k].Sent.Rate())),
					units.BytesSize(float64(packetsStats[k].Received.Rate())),
				})
			}

			print("\033[H\033[2J")
			table.Render()
		}
	}()

	for packet := range p.data {
		if _, exist := packetsStats[packet.Process.Id]; !exist {
			packetsStats[packet.Process.Id] = &Stats{
				Sent:     ratecounter.NewRateCounter(1 * time.Second),
				Received: ratecounter.NewRateCounter(1 * time.Second),
				Packets:  ratecounter.NewRateCounter(1 * time.Second),
			}
		}
		if packet.Process.Sender {
			packetsStats[packet.Process.Id].Sent.Incr(int64(packet.Packet.Size))
		} else {
			packetsStats[packet.Process.Id].Received.Incr(int64(packet.Packet.Size))
		}
		packetsStats[packet.Process.Id].ProcessName = packet.Process.Name
		packetsStats[packet.Process.Id].Packets.Incr(1)
	}
}
