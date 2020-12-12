package process_test

import (
	"net"
	"testing"

	"github.com/maxbaldin/dissertation-project/src/implementation/agent/usecase/process"
	"github.com/stretchr/testify/assert"
)

func TestParseNetTCPRow(t *testing.T) {
	raw := `46: 0401A8C0:01BB 030310AC:1770 01 00000150:00000001  01:00000019 00000001  1000 9 54165785 4 cd1e6040 25 4 27 3 -1`
	row, err := process.ParseNetTCPRow(raw)
	expected := process.NetTCPRow{
		NumberOfEntry: 46,
		Local: process.NetTCPRowAddress{
			IpV4: net.IPv4(192, 168, 1, 4),
			Port: 443,
		},
		Remote: process.NetTCPRowAddress{
			IpV4: net.IPv4(172, 16, 3, 3),
			Port: 6000,
		},
		ConnectionState: 1,
		Queue: process.NetTCPRowQueue{
			Transmission: 336,
			Receive:      1,
		},
		Timer: process.NetTCPRowTimer{
			TimerActive:              1,
			NumOfJiffiesUntilExpires: 25,
		},
		NumOfUnRecoveredRTOTimeouts: 1,
		UID:                         1000,
		UnansweredZeroWinProbes:     9,
		Inode:                       54165785,
		SocketReferenceCount:        4,
		LocationOfSocketInMemory:    "cd1e6040",
		RetransmitTimeout:           25,
		PredictedTickOfSoftClock:    4,
		Ack:                         27,
		SendingCongestionWindow:     3,
		SlowStartSizeThreshold:      -1,
	}
	assert.Equal(t, nil, err)
	assert.Equal(t, expected, row)
}

func TestParseNetTCPRows(t *testing.T) {
	input := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode                                                     
   0: 3500007F:0035 00000000:0000 0A 00000000:00000000 00:00000000 00000000   101        0 103719 1 ffff981d7a8dc000 100 0 0 10 0                    
   1: 0100007F:0277 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 261259 1 ffff981d76e7f800 100 0 0 10 0 `
	out, err := process.ParseNetTCPRows(input)
	assert.Equal(t, nil, err)

	expectedOut := []process.NetTCPRow{
		{
			NumberOfEntry: 0,
			Local: process.NetTCPRowAddress{
				IpV4: net.IPv4(127, 0, 0, 53),
				Port: 53,
			},
			Remote: process.NetTCPRowAddress{
				IpV4: net.IPv4(0, 0, 0, 0),
				Port: 0,
			},
			ConnectionState: 10,
			Queue: process.NetTCPRowQueue{
				Transmission: 0,
				Receive:      0,
			},
			Timer: process.NetTCPRowTimer{
				TimerActive:              0,
				NumOfJiffiesUntilExpires: 0,
			},
			NumOfUnRecoveredRTOTimeouts: 0,
			UID:                         101,
			UnansweredZeroWinProbes:     0,
			Inode:                       103719,
			SocketReferenceCount:        1,
			LocationOfSocketInMemory:    "ffff981d7a8dc000",
			RetransmitTimeout:           100,
			PredictedTickOfSoftClock:    0,
			Ack:                         0,
			SendingCongestionWindow:     10,
			SlowStartSizeThreshold:      0,
		},
		{
			NumberOfEntry: 1,
			Local: process.NetTCPRowAddress{
				IpV4: net.IPv4(127, 0, 0, 1),
				Port: 631,
			},
			Remote: process.NetTCPRowAddress{
				IpV4: net.IPv4(0, 0, 0, 0),
				Port: 0,
			},
			ConnectionState: 10,
			Queue: process.NetTCPRowQueue{
				Transmission: 0,
				Receive:      0,
			},
			Timer: process.NetTCPRowTimer{
				TimerActive:              0,
				NumOfJiffiesUntilExpires: 0,
			},
			NumOfUnRecoveredRTOTimeouts: 0,
			UID:                         0,
			UnansweredZeroWinProbes:     0,
			Inode:                       261259,
			SocketReferenceCount:        1,
			LocationOfSocketInMemory:    "ffff981d76e7f800",
			RetransmitTimeout:           100,
			PredictedTickOfSoftClock:    0,
			Ack:                         0,
			SendingCongestionWindow:     10,
			SlowStartSizeThreshold:      0,
		},
	}

	assert.Equal(t, expectedOut, out)
}
