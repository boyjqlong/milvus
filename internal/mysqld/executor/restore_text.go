package executor

import (
	"bytes"

	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/format"
)

type textRestorer struct{}

func (r textRestorer) getText(node ast.Node) (string, error) {
	buf := bytes.NewBuffer(nil)
	restoreCtx := format.NewRestoreCtx(0, buf)
	if err := node.Restore(restoreCtx); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func newTextRestorer() *textRestorer {
	return &textRestorer{}
}
