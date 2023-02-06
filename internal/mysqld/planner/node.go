package planner

type Node interface {
	GetLocation() NodeLocation
	GetChildren() []Node
}

type baseNode struct {
	location NodeLocation
}

func (n baseNode) GetLocation() NodeLocation {
	return n.location
}

func newBaseNode(location NodeLocation) *baseNode {
	return &baseNode{location: location}
}
