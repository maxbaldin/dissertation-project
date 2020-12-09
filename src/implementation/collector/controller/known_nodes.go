package controller

import (
	"net/http"
	"strings"
)

type KnownNodesRepository interface {
	Get() []string
}

type KnownNodesController struct {
	repository KnownNodesRepository
}

func NewKnownNodesController(repository KnownNodesRepository) *KnownNodesController {
	return &KnownNodesController{
		repository: repository,
	}
}

func (k *KnownNodesController) Handle(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(strings.Join(k.repository.Get(), ",")))
}
