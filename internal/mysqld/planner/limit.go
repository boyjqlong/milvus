package planner

import "fmt"

type NodeLimitClause struct {
	baseNode
	Limit  int
	Offset int
}

func (n *NodeLimitClause) String() string {
	return fmt.Sprintf("Offet: %d, Limit: %d", n.Offset, n.Limit)
}

func (n *NodeLimitClause) GetChildren() []Node {
	return nil
}

func (n *NodeLimitClause) Accept(v Visitor) interface{} {
	return v.VisitLimitClause(n)
}

func NewNodeLimitClause(text string, limit, offset int) *NodeLimitClause {
	return &NodeLimitClause{
		baseNode: newBaseNode(text),
		Limit:    limit,
		Offset:   offset,
	}
}
