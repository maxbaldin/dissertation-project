package collector_test

import (
	"testing"
	"time"

	"github.com/maxbaldin/dissertation-project/src/implementation/agent/integration/collector"
	"github.com/stretchr/testify/assert"
)

func TestNodesRepository_IsKnownNode(t *testing.T) {
	nodesRepo, err := collector.NewNodeRepository("http://127.0.0.1:8080/known_nodes", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 3)
	assert.Equal(t, nodesRepo.IsKnownNode("127"), false)
}
