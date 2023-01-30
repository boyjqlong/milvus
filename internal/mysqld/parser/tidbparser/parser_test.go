package tidbparser

import (
	"fmt"
	"testing"

	"github.com/pingcap/tidb/parser/ast"
)

func getOutputFields(fields *ast.FieldList) ([]string, error) {
	outputs := make([]string, 0, len(fields.Fields))
	for _, field := range fields.Fields {
		if field.WildCard != nil {
			return nil, fmt.Errorf("complex target entry is not supported: %s", field.OriginalText())
		}
		columnNameExpr, ok := field.Expr.(*ast.ColumnNameExpr)
		if !ok {
			return nil, fmt.Errorf("only column as target entry is not supported: %s", field.OriginalText())
		}
		outputs = append(outputs, columnNameExpr.Name.String())
	}
	return outputs, nil
}

func Test_tidbParser_Parse(t *testing.T) {
	/*
		p := tparser.New()

		stmts, warns, err := p.Parse("select * from t where a > 5 and b < 2", "", "")

		stmts, warns, err = p.Parse("select a, b from t where a > 5 and b < 2", "", "")

		buf := bytes.NewBuffer(nil)
		formatter := format.NewRestoreCtx(0, buf)
		err = stmts[0].(*ast.SelectStmt).Where.Restore(formatter)
		fmt.Println(err)
		filter := buf.String()
		fmt.Println(filter)

		outputs, err := getOutputFields(stmts[0].(*ast.SelectStmt).Fields)
		fmt.Println(err)
		fmt.Println(outputs)

		stmts, warns, err = p.Parse("select count(*) from t", "", "")

		fmt.Println(stmts[0].(*ast.SelectStmt).Fields.Fields[0].Expr.(*ast.AggregateFuncExpr).Args[0].Text())

		stmts, warns, err = p.Parse("select count(*) from t where a > 5 and b < 2", "", "")

		if err != nil {
			panic(err)
		}

		for _, warn := range warns {
			fmt.Println(warn)
		}

		fmt.Println(stmts)
	*/
}
