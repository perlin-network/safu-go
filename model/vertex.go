package model

type Vertex struct {
	Address  string
	Children map[string]struct{}
	Parents  map[string]struct{}
}

func NewVertex(address string) *Vertex {
	v := &Vertex{
		Address:  address,
		Parents:  make(map[string]struct{}),
		Children: make(map[string]struct{}),
	}

	return v
}
