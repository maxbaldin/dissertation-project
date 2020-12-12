package entity

import "fmt"

type Edges []Edge

func (n Edges) Exist(id string) bool {
	for _, v := range n {
		if v.Data.Id == id {
			return true
		}
	}
	return false
}

type Edge struct {
	Data struct {
		Id     string `json:"id"`
		Source string `json:"source"`
		Target string `json:"target"`
		Label  string `json:"label"`
		Width  string `json:"width"`
	} `json:"data"`
}

func NewEdge(id, source, target, label string, width int) Edge {
	return Edge{Data: struct {
		Id     string `json:"id"`
		Source string `json:"source"`
		Target string `json:"target"`
		Label  string `json:"label"`
		Width  string `json:"width"`
	}{Id: id, Source: source, Target: target, Label: label, Width: fmt.Sprintf("%dpx", width)}}
}
