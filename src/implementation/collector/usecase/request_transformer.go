package usecase

import (
	"net/http"

	"github.com/maxbaldin/dissertation-project/src/implementation/collector/entity"
)

type RequestTransformer struct {
}

func NewRequestTransformer() *RequestTransformer {
	return &RequestTransformer{}
}

func (rt *RequestTransformer) Transform(r *http.Request) (entity.Traffic, error) {
	return entity.Traffic{}, nil
}
