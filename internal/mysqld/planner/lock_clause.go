package planner

type LockClauseOption int

const (
	LockClauseOption_Unknown LockClauseOption = iota
	LockClauseOption_ForUpdate
	LockClauseOption_LockInShareMode
)

func (o LockClauseOption) String() string {
	switch o {
	case LockClauseOption_ForUpdate:
		return "FOR UPDATE"
	case LockClauseOption_LockInShareMode:
		return "LOCK IN SHARE MODE"
	default:
		return "UNKNOWN"
	}
}

type NodeLockClause struct {
	baseNode
	Option LockClauseOption
}

func (n *NodeLockClause) String() string {
	return n.Option.String()
}

func (n *NodeLockClause) GetChildren() []Node {
	return nil
}

func (n *NodeLockClause) Accept(v Visitor) interface{} {
	return v.VisitLockClause(n)
}

func NewNodeLockClause(text string, option LockClauseOption) *NodeLockClause {
	return &NodeLockClause{
		baseNode: newBaseNode(text),
		Option:   option,
	}
}
