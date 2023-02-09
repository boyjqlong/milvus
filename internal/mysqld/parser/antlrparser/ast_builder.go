package antlrparser

import (
	"fmt"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"

	parsergen "github.com/milvus-io/milvus/internal/mysqld/parser/antlrparser/parser"
	"github.com/milvus-io/milvus/internal/mysqld/planner"
)

type AstBuilder struct {
	parsergen.BaseMySqlParserVisitor
}

type AstBuilderOption func(*AstBuilder)

func (v *AstBuilder) apply(opts ...AstBuilderOption) {
	for _, opt := range opts {
		opt(v)
	}
}

func (v *AstBuilder) VisitRoot(ctx *parsergen.RootContext) interface{} {
	return ctx.SqlStatements().Accept(v)
}

func (v *AstBuilder) VisitSqlStatements(ctx *parsergen.SqlStatementsContext) interface{} {
	allSqlStatementsCtx := ctx.AllSqlStatement()
	sqlStatements := make([]*planner.NodeSqlStatement, 0, len(allSqlStatementsCtx))
	for _, sqlStatementCtx := range allSqlStatementsCtx {
		r := sqlStatementCtx.Accept(v)
		if err := GetError(r); err != nil {
			return err
		}
		n := GetSqlStatement(r)
		if n == nil {
			return fmt.Errorf("failed to parse sql statement: %s", sqlStatementCtx.GetText())
		}
		sqlStatements = append(sqlStatements, n)
	}
	return planner.NewNodeSqlStatements(sqlStatements, ctx.GetText())
}

func (v *AstBuilder) VisitSqlStatement(ctx *parsergen.SqlStatementContext) interface{} {
	dmlStatement := ctx.DmlStatement()
	if dmlStatement == nil {
		return fmt.Errorf("sql statement only support dml statement now: %s", ctx.GetText())
	}
	r := dmlStatement.Accept(v)
	if err := GetError(r); err != nil {
		return err
	}
	n := GetDmlStatement(r)
	if n == nil {
		return fmt.Errorf("failed to parse dml statement: %s", dmlStatement.GetText())
	}
	return planner.NewNodeSqlStatement(ctx.GetText(), planner.WithDmlStatement(n))
}

func (v *AstBuilder) VisitEmptyStatement_(ctx *parsergen.EmptyStatement_Context) interface{} {
	// Should not be visited.
	return nil
}

func (v *AstBuilder) VisitDmlStatement(ctx *parsergen.DmlStatementContext) interface{} {
	selectStatement := ctx.SelectStatement()
	if selectStatement == nil {
		return fmt.Errorf("dml statement only support select statement now: %s", ctx.GetText())
	}
	r := selectStatement.Accept(v)
	if err := GetError(r); err != nil {
		return err
	}
	n := GetSelectStatement(r)
	if n == nil {
		return fmt.Errorf("failed to parse select statement: %s", selectStatement.GetText())
	}
	return planner.NewNodeDmlStatement(ctx.GetText(), planner.WithSelectStatement(n))
}

func (v *AstBuilder) VisitSimpleSelect(ctx *parsergen.SimpleSelectContext) interface{} {
	var opts []planner.NodeSimpleSelectOption

	querySpecificationCtx := ctx.QuerySpecification()
	if querySpecificationCtx == nil {
		return fmt.Errorf("simple select only support query specification now: %s", ctx.GetText())
	}

	r := querySpecificationCtx.Accept(v)
	if err := GetError(r); err != nil {
		return err
	}

	n := GetQuerySpecification(r)
	if n == nil {
		return fmt.Errorf("failed to parse query specification: %s", querySpecificationCtx.GetText())
	}

	opts = append(opts, planner.WithQuery(n))

	lockClauseCtx := ctx.LockClause()
	if lockClauseCtx != nil {
		r := lockClauseCtx.Accept(v)
		if err := GetError(r); err != nil {
			return err
		}

		if n := GetLockClause(r); n != nil {
			opts = append(opts, planner.WithLockClause(n))
		}
	}

	nodeSimpleSelect := planner.NewNodeSimpleSelect(ctx.GetText(), opts...)
	return planner.NewNodeSelectStatement(ctx.GetText(), planner.WithSimpleSelect(nodeSimpleSelect))
}

func (v *AstBuilder) VisitLockClause(ctx *parsergen.LockClauseContext) interface{} {
	//o := planner.LockClauseOption_Unknown
	//text := ctx.GetText()
	//if CaseInsensitiveEqual(text, planner.LockClauseOption_ForUpdate_Str) {
	//	o = planner.LockClauseOption_ForUpdate
	//} else if CaseInsensitiveEqual(text, planner.LockClauseOption_LockInShareMode_Str) {
	//	o = planner.LockClauseOption_LockInShareMode
	//}
	//return planner.NewNodeLockClause(text, o)
	return fmt.Errorf("lock clause is not supported: %s", ctx.GetText())
}

func (v *AstBuilder) VisitQuerySpecification(ctx *parsergen.QuerySpecificationContext) interface{} {
	allSelectSpecCtx := ctx.AllSelectSpec()
	selectSpecs := make([]*planner.NodeSelectSpec, 0, len(allSelectSpecCtx))
	for _, selectSpecCtx := range allSelectSpecCtx {
		r := selectSpecCtx.Accept(v)
		if err := GetError(r); err != nil {
			return err
		}
		if n := GetSelectSpec(r); n != nil {
			selectSpecs = append(selectSpecs, n)
		}
	}

	selectElementsCtx := ctx.SelectElements()
	allSelectElementsCtx := selectElementsCtx.(*parsergen.SelectElementsContext).AllSelectElement()
	selectElements := make([]*planner.NodeSelectElement, 0, len(allSelectElementsCtx))
	if star := selectElementsCtx.GetStar(); star != nil {
		selectElements = append(selectElements, planner.NewNodeSelectElement(star.GetText(), planner.WithStar()))
	}

	for _, selectElementCtx := range allSelectElementsCtx {
		r := selectElementCtx.Accept(v)
		if err := GetError(r); err != nil {
			return err
		}
		if n := GetSelectElement(r); n != nil {
			selectElements = append(selectElements, n)
		}
	}

	var opts []planner.NodeQuerySpecificationOption

	fromCtx := ctx.FromClause()
	if fromCtx != nil {
		r := fromCtx.Accept(v)
		if err := GetError(r); err != nil {
			return err
		}
		if n := GetFromClause(r); n != nil {
			opts = append(opts, planner.WithFrom(n))
		}
	}

	limitCtx := ctx.LimitClause()
	if limitCtx != nil {
		r := limitCtx.Accept(v)
		if err := GetError(r); err != nil {
			return err
		}
		if n := GetLimitClause(r); n != nil {
			opts = append(opts, planner.WithLimit(n))
		}
	}

	return planner.NewNodeQuerySpecification(ctx.GetText(), selectSpecs, selectElements, opts...)
}

func (v *AstBuilder) VisitSelectSpec(ctx *parsergen.SelectSpecContext) interface{} {
	return fmt.Errorf("select spec is not supported: %s", ctx.GetText())
}

func (v *AstBuilder) VisitSelectElements(ctx *parsergen.SelectElementsContext) interface{} {
	// Should not be visited.
	return nil
}

func (v *AstBuilder) VisitSelectColumnElement(ctx *parsergen.SelectColumnElementContext) interface{} {
	var opts []planner.NodeFullColumnNameOption

	asCtx := ctx.AS()
	uidCtx := ctx.Uid()

	if asCtx != nil && uidCtx != nil {
		r := uidCtx.Accept(v)
		if err := GetError(r); err != nil {
			return err
		}
		if alias, ok := r.(string); ok {
			opts = append(opts, planner.FullColumnNameWithAlias(alias))
		}
	}

	if asCtx == nil && uidCtx != nil {
		return fmt.Errorf("invalid alias of column: %s", ctx.GetText())
	}

	r := ctx.FullColumnName().Accept(v)
	if err := GetError(r); err != nil {
		return err
	}
	s, ok := r.(string)
	if !ok {
		return fmt.Errorf("failed to parse full column name: %s", ctx.GetText())
	}
	n := planner.NewNodeFullColumnName(ctx.GetText(), s, opts...)

	return planner.NewNodeSelectElement(ctx.GetText(), planner.WithFullColumnName(n))
}

func (v *AstBuilder) VisitSelectFunctionElement(ctx *parsergen.SelectFunctionElementContext) interface{} {
	var opts []planner.NodeFunctionCallOption

	asCtx := ctx.AS()
	uidCtx := ctx.Uid()

	if asCtx != nil && uidCtx != nil {
		r := uidCtx.Accept(v)
		if err := GetError(r); err != nil {
			return err
		}
		if alias, ok := r.(string); ok {
			opts = append(opts, planner.FunctionCallWithAlias(alias))
		}
	}

	if asCtx == nil && uidCtx != nil {
		return fmt.Errorf("invalid alias of function call: %s", ctx.GetText())
	}

	r := ctx.FunctionCall().Accept(v)
	if err := GetError(r); err != nil {
		return err
	}
	if n := GetAggregateWindowedFunction(r); n != nil {
		opts = append(opts, planner.WithAgg(n))
	}

	n := planner.NewNodeFunctionCall(ctx.GetText(), opts...)

	return planner.NewNodeSelectElement(ctx.GetText(), planner.WithFunctionCall(n))
}

func (v *AstBuilder) VisitFromClause(ctx *parsergen.FromClauseContext) interface{} {
	text := ctx.GetText()

	tableSourcesCtx := ctx.TableSources().(*parsergen.TableSourcesContext)
	allTableSourcesCtx := tableSourcesCtx.AllTableSource()
	tableSources := make([]*planner.NodeTableSource, 0, len(allTableSourcesCtx))
	for _, tableSourceCtx := range allTableSourcesCtx {
		r := tableSourceCtx.Accept(v)
		if err := GetError(r); err != nil {
			return err
		}
		if n := GetTableSource(r); n != nil {
			tableSources = append(tableSources, n)
		}
	}

	var opts []planner.NodeFromClauseOption

	whereCtx := ctx.WHERE()
	if whereCtx != nil {
		r := whereCtx.Accept(v)
		if err := GetError(r); err != nil {
			return err
		}
		if n := GetExpression(r); n != nil {
			opts = append(opts, planner.WithWhere(n))
		}
	}

	return planner.NewNodeFromClause(text, tableSources, opts...)
}

func (v *AstBuilder) VisitTableSources(ctx *parsergen.TableSourcesContext) interface{} {
	// Should not be visited.
	return nil
}

func (v *AstBuilder) VisitTableSourceBase(ctx *parsergen.TableSourceBaseContext) interface{} {
	text := ctx.GetText()

	tableNameCtx := ctx.TableName()
	r := tableNameCtx.Accept(v)

	if err := GetError(r); err != nil {
		return err
	}

	s, ok := r.(string)
	if !ok {
		return fmt.Errorf("not supported table source: %s", text)
	}

	return planner.NewNodeTableSource(text, planner.WithTableName(s))
}

func (v *AstBuilder) VisitLimitClause(ctx *parsergen.LimitClauseContext) interface{} {
	text := ctx.GetText()
	offset, limit := 0, 0

	getInt := func(tree antlr.ParseTree) (int, error) {
		r := tree.Accept(v)
		if err := GetError(r); err != nil {
			return 0, err
		}
		i, ok := r.(int)
		if !ok {
			return 0, fmt.Errorf("failed to parse to int: %s", tree.GetText())
		}
		return i, nil
	}

	offsetCtx := ctx.GetOffset()
	if offsetCtx != nil {
		i, err := getInt(offsetCtx)
		if err != nil {
			return err
		}
		offset = i
	}

	limitCtx := ctx.GetLimit()
	i, err := getInt(limitCtx)
	if err != nil {
		return err
	}
	limit = i

	return planner.NewNodeLimitClause(text, limit, offset)
}

func (v *AstBuilder) VisitLimitClauseAtom(ctx *parsergen.LimitClauseAtomContext) interface{} {
	text := ctx.GetText()

	decimalLiteralCtx := ctx.DecimalLiteral()
	if decimalLiteralCtx == nil {
		return fmt.Errorf("limit only support literal now: %s", text)
	}

	return decimalLiteralCtx.Accept(v)
}

func (v *AstBuilder) VisitFullId(ctx *parsergen.FullIdContext) interface{} {
	return ctx.GetText()
}

func (v *AstBuilder) VisitTableName(ctx *parsergen.TableNameContext) interface{} {
	return ctx.GetText()
}

func (v *AstBuilder) VisitFullColumnName(ctx *parsergen.FullColumnNameContext) interface{} {
	return ctx.GetText()
}

func (v *AstBuilder) VisitCollationName(ctx *parsergen.CollationNameContext) interface{} {
	return ctx.GetText()
}

func (v *AstBuilder) VisitUid(ctx *parsergen.UidContext) interface{} {
	return ctx.GetText()
}

func (v *AstBuilder) VisitSimpleId(ctx *parsergen.SimpleIdContext) interface{} {
	return ctx.GetText()
}

func (v *AstBuilder) VisitDottedId(ctx *parsergen.DottedIdContext) interface{} {
	return ctx.GetText()
}

func (v *AstBuilder) VisitDecimalLiteral(ctx *parsergen.DecimalLiteralContext) interface{} {
	if decimalCtx := ctx.DECIMAL_LITERAL(); decimalCtx != nil {
		i, err := strconv.Atoi(decimalCtx.GetText())
		if err != nil {
			return err
		}
		return int64(i)
	}

	if zeroCtx := ctx.ZERO_DECIMAL(); zeroCtx != nil {
		return int64(0)
	}

	if oneCtx := ctx.ONE_DECIMAL(); oneCtx != nil {
		return int64(1)
	}

	if twoCtx := ctx.TWO_DECIMAL(); twoCtx != nil {
		return int64(2)
	}

	if realCtx := ctx.REAL_LITERAL(); realCtx != nil {
		f, err := processRealCtx(realCtx)
		if err != nil {
			return err
		}
		return f
	}

	return nil
}

func processRealCtx(n antlr.TerminalNode) (float64, error) {
	f, err := strconv.ParseFloat(n.GetText(), 64)
	if err != nil {
		return 0.0, err
	}
	return f, nil
}

func (v *AstBuilder) VisitStringLiteral(ctx *parsergen.StringLiteralContext) interface{} {
	return ctx.GetText()
}

func (v *AstBuilder) VisitBooleanLiteral(ctx *parsergen.BooleanLiteralContext) interface{} {
	if ctx.TRUE() != nil {
		return true
	}
	return false
}

func (v *AstBuilder) VisitHexadecimalLiteral(ctx *parsergen.HexadecimalLiteralContext) interface{} {
	text := ctx.GetText()
	i, err := strconv.ParseInt(text, 16, 64)
	if err != nil {
		return err
	}
	return i
}

func (v *AstBuilder) VisitConstant(ctx *parsergen.ConstantContext) interface{} {
	if stringLiteralCtx := ctx.StringLiteral(); stringLiteralCtx != nil {
		return stringLiteralCtx.Accept(v)
	}

	if decimalLiteralCtx := ctx.DecimalLiteral(); decimalLiteralCtx != nil {
		decimal := decimalLiteralCtx.Accept(v)
		if err := GetError(decimal); err != nil {
			return err
		}
		if ctx.MINUS() != nil {
			switch realType := decimal.(type) {
			case int:
				return -realType
			case int32:
				return -realType
			case int64:
				return -realType
			case float32:
				return -realType
			case float64:
				return -realType
			default:
				return fmt.Errorf("invalid data type: %s", decimalLiteralCtx.GetText())
			}
		}
		return decimal
	}

	if hexadecimalLiteralCtx := ctx.HexadecimalLiteral(); hexadecimalLiteralCtx != nil {
		return hexadecimalLiteralCtx.Accept(v)
	}

	if realCtx := ctx.REAL_LITERAL(); realCtx != nil {
		f, err := processRealCtx(realCtx)
		if err != nil {
			return err
		}
		return f
	}

	return nil
}

func (v *AstBuilder) VisitExpressions(ctx *parsergen.ExpressionsContext) interface{} {
	text := ctx.GetText()
	expressionsCtx := ctx.AllExpression()
	expressions := make([]*planner.NodeExpression, 0, len(expressionsCtx))
	for _, expressionCtx := range expressionsCtx {
		n, err := v.getExpression(expressionCtx)
		if err != nil {
			return err
		}
		if n == nil {
			return fmt.Errorf("failed to parse expression: %s", expressionCtx.GetText())
		}
		expressions = append(expressions, n)
	}
	return planner.NewNodeExpressions(text, expressions)
}

func (v *AstBuilder) VisitAggregateFunctionCall(ctx *parsergen.AggregateFunctionCallContext) interface{} {
	return ctx.AggregateWindowedFunction().Accept(v)
}

func (v *AstBuilder) VisitAggregateWindowedFunction(ctx *parsergen.AggregateWindowedFunctionContext) interface{} {
	text := ctx.GetText()
	return planner.NewNodeCount(text)
}

func (v *AstBuilder) getPredicate(tree antlr.ParseTree) (*planner.NodePredicate, error) {
	r := tree.Accept(v)
	if err := GetError(r); err != nil {
		return nil, err
	}
	return GetPredicate(r), nil
}

func (v *AstBuilder) VisitIsExpression(ctx *parsergen.IsExpressionContext) interface{} {
	text := ctx.GetText()

	testValue := planner.TestValue_Unknown
	if ctx.TRUE() != nil {
		testValue = planner.TestValue_True
	}
	if ctx.FALSE() != nil {
		testValue = planner.TestValue_False
	}

	op := planner.IsOperator_Is
	if ctx.NOT() != nil {
		op = planner.IsOperator_IsNot
	}

	n, err := v.getPredicate(ctx.Predicate())

	if err != nil {
		return err
	}

	if n == nil {
		return fmt.Errorf("failed to parse [is expression]: %s", text)
	}

	isExpr := planner.NewNodeIsExpression(text, n, testValue, op)
	return planner.NewNodeExpression(text, planner.WithIsExpr(isExpr))
}

func (v *AstBuilder) getExpression(tree antlr.ParseTree) (*planner.NodeExpression, error) {
	r := tree.Accept(v)
	if err := GetError(r); err != nil {
		return nil, err
	}
	return GetExpression(r), nil
}

func (v *AstBuilder) VisitNotExpression(ctx *parsergen.NotExpressionContext) interface{} {
	text := ctx.GetText()
	expressionCtx := ctx.Expression()
	n, err := v.getExpression(expressionCtx)
	if err != nil {
		return err
	}
	if n == nil {
		return fmt.Errorf("failed to parse [not expression]: %s", text)
	}
	notExpr := planner.NewNodeNotExpression(text, n)
	return planner.NewNodeExpression(text, planner.WithNotExpr(notExpr))
}

func (v *AstBuilder) VisitLogicalExpression(ctx *parsergen.LogicalExpressionContext) interface{} {
	text := ctx.GetText()
	op := ctx.LogicalOperator().Accept(v).(planner.LogicalOperator)
	left, right := ctx.Expression(0), ctx.Expression(1)

	leftExpr, err := v.getExpression(left)
	if err != nil {
		return err
	}

	if leftExpr == nil {
		return fmt.Errorf("failed to parse left expr: %s", left.GetText())
	}

	rightExpr, err := v.getExpression(right)
	if err != nil {
		return err
	}

	if rightExpr == nil {
		return fmt.Errorf("failed to parse right expr: %s", right.GetText())
	}

	logicalExpr := planner.NewNodeLogicalExpression(text, leftExpr, rightExpr, op)
	return planner.NewNodeExpression(text, planner.WithLogicalExpr(logicalExpr))
}

func (v *AstBuilder) VisitPredicateExpression(ctx *parsergen.PredicateExpressionContext) interface{} {
	text := ctx.GetText()
	predicateCtx := ctx.Predicate()
	predicate, err := v.getPredicate(predicateCtx)
	if err != nil {
		return err
	}
	if predicate == nil {
		return fmt.Errorf("failed to parse predicate expression: %s", predicateCtx.GetText())
	}
	return planner.NewNodeExpression(text, planner.WithPredicate(predicate))
}

func (v *AstBuilder) getExpressionAtom(tree antlr.ParseTree) (*planner.NodeExpressionAtom, error) {
	r := tree.Accept(v)
	if err := GetError(r); err != nil {
		return nil, err
	}
	return GetExpressionAtom(r), nil
}

func (v *AstBuilder) VisitExpressionAtomPredicate(ctx *parsergen.ExpressionAtomPredicateContext) interface{} {
	text := ctx.GetText()
	if ctx.LOCAL_ID() != nil || ctx.VAR_ASSIGN() != nil {
		return fmt.Errorf("assignment not supported: %s", text)
	}
	n, err := v.getExpressionAtom(ctx.ExpressionAtom())
	if err != nil {
		return err
	}
	if n == nil {
		return fmt.Errorf("failed to parse expression atom predicate: %s", text)
	}
	expr := planner.NewNodeExpressionAtomPredicate(text, n)
	return planner.NewNodePredicate(text, planner.WithNodeExpressionAtomPredicate(expr))
}

func (v *AstBuilder) VisitBinaryComparisonPredicate(ctx *parsergen.BinaryComparisonPredicateContext) interface{} {
	text := ctx.GetText()

	op := ctx.ComparisonOperator().Accept(v).(planner.ComparisonOperator)

	leftCtx, rightCtx := ctx.GetLeft(), ctx.GetRight()

	left, err := v.getPredicate(leftCtx)
	if err != nil {
		return err
	}
	if left == nil {
		return fmt.Errorf("failed to parse left predicate: %s", leftCtx.GetText())
	}

	right, err := v.getPredicate(rightCtx)
	if err != nil {
		return err
	}
	if right == nil {
		return fmt.Errorf("failed to parse right predicate: %s", rightCtx.GetText())
	}

	expr := planner.NewNodeBinaryComparisonPredicate(text, left, right, op)
	return planner.NewNodePredicate(text, planner.WithNodeBinaryComparisonPredicate(expr))
}

func (v *AstBuilder) VisitInPredicate(ctx *parsergen.InPredicateContext) interface{} {
	text := ctx.GetText()

	op := planner.InOperator_In
	if ctx.NOT() != nil {
		op = planner.InOperator_NotIn
	}

	predicateCtx := ctx.Predicate()
	n, err := v.getPredicate(predicateCtx)
	if err != nil {
		return err
	}
	if n == nil {
		return fmt.Errorf("failed to parse left in predicate: %s", predicateCtx.GetText())
	}

	expressionsCtx := ctx.Expressions()
	r := expressionsCtx.Accept(v)
	if err := GetError(r); err != nil {
		return err
	}
	exprs := GetExpressions(r)
	if exprs == nil {
		return fmt.Errorf("failed to parse in expressions: %s", expressionsCtx.GetText())
	}

	expr := planner.NewNodeInPredicate(text, n, exprs, op)
	return planner.NewNodePredicate(text, planner.WithInPredicate(expr))
}

func (v *AstBuilder) VisitConstantExpressionAtom(ctx *parsergen.ConstantExpressionAtomContext) interface{} {
	text := ctx.GetText()
	constantCtx := ctx.Constant().(*parsergen.ConstantContext)
	r := constantCtx.Accept(v)
	if err := GetError(r); err != nil {
		return err
	}
	var opts []planner.NodeConstantOption
	switch t := r.(type) {
	case int64:
		opts = append(opts, planner.WithDecimalLiteral(t))
	case string:
		opts = append(opts, planner.WithStringLiteral(t))
	case float64:
		opts = append(opts, planner.WithRealLiteral(t))
	case bool:
		opts = append(opts, planner.WithBooleanLiteral(t))
	}
	c := planner.NewNodeConstant(constantCtx.GetText(), opts...)
	return planner.NewNodeExpressionAtom(text, planner.ExpressionAtomWithConstant(c))
}

func (v *AstBuilder) VisitFullColumnNameExpressionAtom(ctx *parsergen.FullColumnNameExpressionAtomContext) interface{} {
	text := ctx.GetText()
	fullColumnNameCtx := ctx.FullColumnName()
	r := fullColumnNameCtx.Accept(v)
	if err := GetError(r); err != nil {
		return err
	}
	n := GetFullColumnName(r)
	return planner.NewNodeExpressionAtom(text, planner.ExpressionAtomWithFullColumnName(n))
}

func (v *AstBuilder) VisitUnaryExpressionAtom(ctx *parsergen.UnaryExpressionAtomContext) interface{} {
	text := ctx.GetText()
	op := ctx.UnaryOperator().Accept(v).(planner.UnaryOperator)
	expressionAtomCtx := ctx.ExpressionAtom()
	r := expressionAtomCtx.Accept(v)
	if err := GetError(r); err != nil {
		return err
	}
	n := GetExpressionAtom(r)
	expr := planner.NewNodeUnaryExpressionAtom(expressionAtomCtx.GetText(), n, op)
	return planner.NewNodeExpressionAtom(text, planner.ExpressionAtomWithUnaryExpr(expr))
}

func (v *AstBuilder) VisitNestedExpressionAtom(ctx *parsergen.NestedExpressionAtomContext) interface{} {
	text := ctx.GetText()
	expressionsCtx := ctx.AllExpression()
	expressions := make([]*planner.NodeExpression, 0, len(expressionsCtx))
	for _, expressionCtx := range expressionsCtx {
		n, err := v.getExpression(expressionCtx)
		if err != nil {
			return err
		}
		expressions = append(expressions, n)
	}
	return planner.NewNodeNestedExpressionAtom(text, expressions)
}

func (v *AstBuilder) VisitUnaryOperator(ctx *parsergen.UnaryOperatorContext) interface{} {
	op := planner.UnaryOperator_Unknown
	if ctx.EXCLAMATION_SYMBOL() != nil {
		op = planner.UnaryOperator_ExclamationSymbol
	}
	if ctx.BIT_NOT_OP() != nil {
		op = planner.UnaryOperator_Tilde
	}
	if ctx.PLUS() != nil {
		op = planner.UnaryOperator_Positive
	}
	if ctx.MINUS() != nil {
		op = planner.UnaryOperator_Negative
	}
	if ctx.NOT() != nil {
		op = planner.UnaryOperator_Not
	}
	return op
}

func (v *AstBuilder) VisitComparisonOperator(ctx *parsergen.ComparisonOperatorContext) interface{} {
	op := planner.ComparisonOperator_Unknown
	if ctx.EXCLAMATION_SYMBOL() != nil {
		op = planner.ComparisonOperator_NotEqual
	}
	if ctx.EQUAL_SYMBOL() != nil {
		op |= planner.ComparisonOperator_Equal
	}
	if ctx.LESS_SYMBOL() != nil {
		op |= planner.ComparisonOperator_LessThan
	}
	if ctx.GREATER_SYMBOL() != nil {
		op |= planner.ComparisonOperator_GreaterThan
	}
	return op
}

func (v *AstBuilder) VisitLogicalOperator(ctx *parsergen.LogicalOperatorContext) interface{} {
	op := planner.LogicalOperator_Unknown
	if ctx.AND() != nil {
		op = planner.LogicalOperator_And
	}
	if len(ctx.AllBIT_AND_OP()) != 0 {
		op = planner.LogicalOperator_And
	}
	if ctx.OR() != nil {
		op = planner.LogicalOperator_Or
	}
	if len(ctx.AllBIT_OR_OP()) != 0 {
		op = planner.LogicalOperator_Or
	}
	return op
}

func (v *AstBuilder) Build(ctx parsergen.IRootContext) (*planner.NodeSqlStatements, error) {
	n := ctx.Accept(v)
	if err := GetError(n); err != nil {
		return nil, err
	}
	if statements := GetSqlStatements(n); statements != nil {
		return statements, nil
	}
	return nil, fmt.Errorf("failed to parse ast tree: %s", ctx.GetText())
}

func NewAstBuilder(opts ...AstBuilderOption) parsergen.MySqlParserVisitor {
	v := &AstBuilder{}
	v.apply(opts...)
	return v
}
