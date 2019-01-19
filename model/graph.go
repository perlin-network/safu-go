package model

import "fmt"

type DAG struct {
	vertices map[string]*Vertex
}

func NewDAG() *DAG {
	d := &DAG{
		vertices: make(map[string]*Vertex),
	}

	return d
}

func (d *DAG) AddVertex(v *Vertex) {
	d.vertices[v.Address] = v
}

func (d *DAG) DeleteVertex(vertex *Vertex) error {
	if _, ok := d.vertices[vertex.Address]; !ok {
		return fmt.Errorf("vertex with ID %v does not exist", vertex.Address)
	}

	delete(d.vertices, vertex.Address)

	return nil
}

//func (d *DAG) AddEdge(parent *Vertex, child *Vertex) error {
//	if _, ok := d.vertices[parent.Address]; !ok {
//		return fmt.Errorf("vertex %v does not exist", parent.Address)
//	}
//
//	if _, ok := d.vertices[child.Address]; !ok {
//		return fmt.Errorf("vertex ID %v does not exist", child.Address)
//	}
//
//	if _, ok := parent.Children[child.Address]; ok {
//		return fmt.Errorf("edge (%v,%v) already exists", parent.Address, child.Address)
//	}
//
//	// Add edge.
//	parent.Children[child.Address] = struct{}{}
//	child.Parents[parent.Address] = struct{}{}
//
//	return nil
//}