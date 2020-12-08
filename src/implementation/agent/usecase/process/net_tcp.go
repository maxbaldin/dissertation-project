package process

import (
	"encoding/hex"
	"errors"
	"log"
	"net"
	"strconv"
	"strings"
)

type (
	// https://www.kernel.org/doc/Documentation/networking/proc_net_tcp.txt
	NetTCPRow struct {
		NumberOfEntry               int
		Local                       NetTCPRowAddress
		Remote                      NetTCPRowAddress
		ConnectionState             int
		Queue                       NetTCPRowQueue
		Timer                       NetTCPRowTimer
		NumOfUnRecoveredRTOTimeouts int
		UID                         int
		UnansweredZeroWinProbes     int
		Inode                       int
		SocketReferenceCount        int
		LocationOfSocketInMemory    string
		RetransmitTimeout           int
		PredictedTickOfSoftClock    int
		Ack                         int
		SendingCongestionWindow     int
		SlowStartSizeThreshold      int
	}
	NetTCPRowTimer struct {
		TimerActive              int
		NumOfJiffiesUntilExpires int
	}
	NetTCPRowAddress struct {
		IpV4 net.IP
		Port int
	}
	NetTCPRowQueue struct {
		Transmission int
		Receive      int
	}
)

var (
	ErrNetTcpInvalidRawRow    = errors.New("net/tcp: invalid raw row")
	ErrNetTcpUnableToParseRow = errors.New("net/tcp: unable to parse row")
)

func ParseNetTCPRows(raw string) (rows []NetTCPRow, err error) {
	parts := strings.Split(raw, "\n")
	data := parts[1:]
	for _, row := range data {
		trimmedRow := strings.TrimSpace(row)
		if len(trimmedRow) == 0 {
			continue
		}
		parsed, err := ParseNetTCPRow(trimmedRow)
		if err == ErrNetTcpInvalidRawRow {
			continue
		}
		if err != nil {
			return rows, err
		}
		rows = append(rows, parsed)
	}
	return
}

func ParseNetTCPRow(raw string) (row NetTCPRow, err error) {
	parts := cleanSlice(strings.Split(raw, " "))
	if len(parts) != 17 {
		return row, ErrNetTcpInvalidRawRow
	}

	row.NumberOfEntry, err = strconv.Atoi(strings.TrimRight(parts[0], ":"))
	if err != nil {
		log.Println("Wrong NumberOfEntry", parts[0])
		return row, ErrNetTcpUnableToParseRow
	}

	row.Local, err = parseAddress(parts[1])
	if err != nil {
		log.Println("Wrong Local", parts[1])
		return row, ErrNetTcpUnableToParseRow
	}

	row.Remote, err = parseAddress(parts[2])
	if err != nil {
		log.Println("Wrong Remote", parts[2])
		return row, ErrNetTcpUnableToParseRow
	}

	row.ConnectionState, err = hex2int(parts[3])
	if err != nil {
		log.Println("Wrong ConnectionState", parts[3])
		return row, ErrNetTcpUnableToParseRow
	}

	row.Queue, err = parseQueue(parts[4])
	if err != nil {
		log.Println("Wrong Queue", parts[4])
		return row, ErrNetTcpUnableToParseRow
	}

	row.Timer, err = parseTimer(parts[5])
	if err != nil {
		log.Println("Wrong Timer", parts[5])
		return row, ErrNetTcpUnableToParseRow
	}

	row.NumOfUnRecoveredRTOTimeouts, err = strconv.Atoi(parts[6])
	if err != nil {
		log.Println("Wrong NumOfUnRecoveredRTOTimeouts", parts[6])
		return row, ErrNetTcpUnableToParseRow
	}

	row.UID, err = strconv.Atoi(parts[7])
	if err != nil {
		log.Println("Wrong UID", parts[7])
		return row, ErrNetTcpUnableToParseRow
	}

	row.UnansweredZeroWinProbes, err = strconv.Atoi(parts[8])
	if err != nil {
		log.Println("Wrong UnansweredZeroWinProbes", parts[8])
		return row, ErrNetTcpUnableToParseRow
	}

	row.Inode, err = strconv.Atoi(parts[9])
	if err != nil {
		log.Println("Wrong Inode", parts[9])
		return row, ErrNetTcpUnableToParseRow
	}

	row.SocketReferenceCount, err = strconv.Atoi(parts[10])
	if err != nil {
		log.Println("Wrong SocketReferenceCount", parts[10])
		return row, ErrNetTcpUnableToParseRow
	}

	row.LocationOfSocketInMemory = parts[11]

	row.RetransmitTimeout, err = strconv.Atoi(parts[12])
	if err != nil {
		log.Println("Wrong RetransmitTimeout", parts[12])
		return row, ErrNetTcpUnableToParseRow
	}

	row.PredictedTickOfSoftClock, err = strconv.Atoi(parts[13])
	if err != nil {
		log.Println("Wrong PredictedTickOfSoftClock", parts[13])
		return row, ErrNetTcpUnableToParseRow
	}

	row.Ack, err = strconv.Atoi(parts[14])
	if err != nil {
		log.Println("Wrong Ack", parts[14])
		return row, ErrNetTcpUnableToParseRow
	}

	row.SendingCongestionWindow, err = strconv.Atoi(parts[15])
	if err != nil {
		log.Println("Wrong SendingCongestionWindow", parts[15])
		return row, ErrNetTcpUnableToParseRow
	}

	row.SlowStartSizeThreshold, err = strconv.Atoi(parts[16])
	if err != nil {
		log.Println("Wrong SlowStartSizeThreshold", parts[16])
		return row, ErrNetTcpUnableToParseRow
	}

	return
}

func parseAddress(raw string) (addr NetTCPRowAddress, err error) {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return addr, ErrNetTcpUnableToParseRow
	}
	ip, err := hex.DecodeString(parts[0])
	if err != nil || len(ip) != 4 {
		return addr, ErrNetTcpUnableToParseRow
	}
	addr.IpV4 = net.IPv4(ip[3], ip[2], ip[1], ip[0])

	addr.Port, err = hex2int(parts[1])
	if err != nil {
		return addr, ErrNetTcpUnableToParseRow
	}

	return
}

func parseQueue(raw string) (queue NetTCPRowQueue, err error) {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return queue, ErrNetTcpUnableToParseRow
	}
	queue.Transmission, err = hex2int(parts[0])
	if err != nil {
		return queue, ErrNetTcpUnableToParseRow
	}
	queue.Receive, err = hex2int(parts[1])
	if err != nil {
		return queue, ErrNetTcpUnableToParseRow
	}
	return
}

func parseTimer(raw string) (timer NetTCPRowTimer, err error) {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return timer, ErrNetTcpUnableToParseRow
	}
	timer.TimerActive, err = hex2int(parts[0])
	if err != nil {
		return timer, ErrNetTcpUnableToParseRow
	}
	timer.NumOfJiffiesUntilExpires, err = hex2int(parts[1])
	if err != nil {
		return timer, ErrNetTcpUnableToParseRow
	}
	return
}

func hex2int(hexStr string) (int, error) {
	// remove 0x suffix if found in the input string
	cleaned := strings.Replace(hexStr, "0x", "", -1)
	val, err := strconv.ParseInt(cleaned, 16, 64)
	return int(val), err
}

func cleanSlice(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
