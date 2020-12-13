package usecase

import (
	"fmt"

	"github.com/docker/go-units"
	"github.com/maxbaldin/dissertation-project/src/implementation/ui/entity"
	"github.com/maxbaldin/dissertation-project/src/implementation/ui/integration/mysql"
)

const WeightClustersCnr = 6

type NodeTransformer struct {
}

func NewNodeTransformer() *NodeTransformer {
	return &NodeTransformer{}
}

func (t *NodeTransformer) Transform(in entity.GraphDataElements) (resp entity.GraphResponse, err error) {
	response := entity.GraphResponse{
		Nodes: nil,
		Edges: nil,
	}

	for _, v := range in {
		// nodes
		if !response.Nodes.Exist(v.SourceHost) {
			node := entity.NewNode(v.SourceHost, "")
			response.Nodes = append(response.Nodes, node)
		}
		if !response.Nodes.Exist(v.TargetHost) {
			node := entity.NewNode(v.TargetHost, "")
			response.Nodes = append(response.Nodes, node)
		}
		if !response.Nodes.Exist(v.SourceServiceID()) {
			node := entity.NewNode(v.SourceServiceID(), v.SourceHost)
			response.Nodes = append(response.Nodes, node)
		}
		if !response.Nodes.Exist(v.TargetServiceID()) {
			node := entity.NewNode(v.TargetServiceID(), v.TargetHost)
			response.Nodes = append(response.Nodes, node)
		}

		// edges
		if !response.Edges.Exist(v.EdgeId()) {
			weight := int(v.SizeBytes / (in.MaxSize() / WeightClustersCnr))
			if weight < 1 {
				weight = 1
			}
			composedLabel := fmt.Sprintf(
				"%s / %.2fMbps",
				units.BytesSize(float64(v.SizeBytes)),
				(float32(v.SizeBytes)/(60*60*mysql.LastHoursCnt)/1024/1024)*8,
			)
			edge := entity.NewEdge(
				v.EdgeId(),
				v.SourceServiceID(),
				v.TargetServiceID(),
				composedLabel,
				weight,
			)
			response.Edges = append(response.Edges, edge)
		}
	}
	return response, nil
}
