package planner

import "github.com/pingcap/tidb/parser/ast"

type LogicalPlan struct {
	Stmt ast.StmtNode
}

type PhysicalPlan struct {
	Stmt ast.StmtNode
}
