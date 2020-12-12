package repository

import "github.com/maxbaldin/dissertation-project/src/implementation/ui/entity"

type NodesDb interface {
	GraphData() (entity.GraphDataElements, error)
}

type GraphDataTransformer interface {
	Transform(elements entity.GraphDataElements) (entity.GraphResponse, error)
}

type NodesRepository struct {
	db          NodesDb
	transformer GraphDataTransformer
}

func NewNodesRepository(db NodesDb, transformer GraphDataTransformer) *NodesRepository {
	return &NodesRepository{
		db:          db,
		transformer: transformer,
	}
}

func (n *NodesRepository) FindAll() (resp entity.GraphResponse, err error) {
	data, err := n.db.GraphData()
	if err != nil {
		return resp, nil
	}
	return n.transformer.Transform(data)
}
