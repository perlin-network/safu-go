package model

import (
	"bytes"
)

type Vertex struct {
	Address  string              `json:"address"`
	Taint    int                 `json:"taint"`
	Children map[string]struct{} `json:"children"`
	Parents  map[string]struct{} `json:"parents"`
}

func NewVertex(address string) *Vertex {
	v := &Vertex{
		Address:  address,
		Parents:  make(map[string]struct{}),
		Children: make(map[string]struct{}),
	}

	return v
}

func (v *Vertex) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("address: " + v.Address)

	for k := range v.Children {
		buffer.WriteString("		children: " + k)
	}
	return buffer.String()
}
