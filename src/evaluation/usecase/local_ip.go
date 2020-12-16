package usecase

import (
	"net"
	"strings"
)

func ReplaceLocalhostWithOutboundIP(in string) string {
	if strings.Contains(in, "host") {
		in = strings.Replace(in, "host", GetOutboundIP().String(), 1)
	}
	return in
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
