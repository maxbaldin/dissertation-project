package controller

import (
	"fmt"
	"net/http"
)

type Index struct {
}

func NewIndex() *Index {
	return &Index{}
}

func (i *Index) Handle(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
