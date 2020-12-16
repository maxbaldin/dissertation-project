package usecase

import (
	"errors"
	"fmt"

	"github.com/docker/go-units"
	"github.com/maxbaldin/dissertation-project/src/implementation/ui/entity"
)

const WeightClustersCnr = 6

var ErrNoData = errors.New("no data")

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
				(float32(v.SizeBytesPerMinute)/1024/1024/60)*8,
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

	if len(response.Edges) == 0 && len(response.Nodes) == 0 {
		return resp, ErrNoData
	}

	return response, nil
}
