package entity

type Nodes []Node

func (n Nodes) Exist(id string) bool {
	for _, v := range n {
		if v.Data.Id == id {
			return true
		}
	}
	return false
}

type Node struct {
	Data struct {
		Id     string `json:"id"`
		Parent string `json:"parent,omitempty"`
	} `json:"data"`
}

func NewNode(id, parent string) Node {
	return Node{Data: struct {
		Id     string `json:"id"`
		Parent string `json:"parent,omitempty"`
	}{Id: id, Parent: parent}}
}
