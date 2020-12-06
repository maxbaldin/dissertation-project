package process

import "io/ioutil"

type (
	NetTCPRowsIndex map[string]map[int]NetTCPRow
)

func NewNetTCPIndex() (index NetTCPRowsIndex, err error) {
	netTcpOut, err := ioutil.ReadFile("/proc/net/tcp")
	if err != nil {
		return
	}
	connections, err := ParseNetTCPRows(string(netTcpOut))
	if err != nil {
		return
	}
	idx := NetTCPRowsIndex{}
	idx.createSourceIdx(connections)
	return idx, nil
}

func (v NetTCPRowsIndex) createSourceIdx(rows []NetTCPRow) {
	for _, row := range rows {
		localIp := row.Local.IpV4.String()
		localPort := row.Local.Port
		if _, exist := v[row.Local.IpV4.String()]; !exist {
			v[localIp] = map[int]NetTCPRow{}
		}
		v[localIp][localPort] = row
	}
}

func (v NetTCPRowsIndex) LookupSource(ip string, port int) (response NetTCPRow, exist bool) {
	if ipMap, ipMapExist := v[ip]; ipMapExist {
		if portMap, portMapExist := ipMap[port]; portMapExist {
			return portMap, true
		}
	}
	return
}
