package tidbparser

import (
	"fmt"

	"github.com/milvus-io/milvus/internal/mysqld/parser"
	"github.com/milvus-io/milvus/internal/mysqld/planner"
	tparser "github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/types/parser_driver" // used for tidb parser.
)

type tidbParser struct{}

func (w tidbParser) Parse(sql string, opts ...parser.Option) (*planner.LogicalPlan, []error, error) {
	cfg := parser.DefaultParserConfig()
	cfg.Apply(opts...)

	p := parserPool.Get().(*tparser.Parser)
	defer parserPool.Put(p)

	p.SetSQLMode(cfg.SQLMode)
	p.SetParserConfig(cfg.ToTidbParserConfig())

	stmts, warns, err := p.Parse(sql, cfg.Charset, cfg.Collation)
	if err != nil {
		return nil, nil, err
	}

	tmp := make([]ast.StmtNode, len(stmts))
	copy(tmp, stmts)

	if len(stmts) != 1 {
		return nil, warns, fmt.Errorf("only one query is allowed: %s", sql)
	}

	return &planner.LogicalPlan{Stmt: tmp[0]}, warns, nil
}

func NewTidbParser() parser.Parser {
	return &tidbParser{}
}
