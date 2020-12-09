package utils_test

import (
	"net"
	"testing"

	"github.com/maxbaldin/dissertation-project/src/implementation/collector/usecase/utils"
	"github.com/stretchr/testify/assert"
)

func TestIp2int(t *testing.T) {
	ip := net.ParseIP("0.0.0.0")
	assert.Equal(t, uint32(0), utils.Ip2int(ip))

	ip = net.ParseIP("172.19.0.2")
	assert.Equal(t, uint32(2886926338), utils.Ip2int(ip))

}
