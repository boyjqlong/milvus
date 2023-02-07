package antlrparser

import (
	parsergen "github.com/milvus-io/milvus/internal/mysqld/parser/antlrparser/parser"
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
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitSqlStatements(ctx *parsergen.SqlStatementsContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitSqlStatement(ctx *parsergen.SqlStatementContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitEmptyStatement_(ctx *parsergen.EmptyStatement_Context) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitDmlStatement(ctx *parsergen.DmlStatementContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitSimpleSelect(ctx *parsergen.SimpleSelectContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitParenthesisSelect(ctx *parsergen.ParenthesisSelectContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitLockClause(ctx *parsergen.LockClauseContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitTableSources(ctx *parsergen.TableSourcesContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitTableSourceBase(ctx *parsergen.TableSourceBaseContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitQueryExpression(ctx *parsergen.QueryExpressionContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitQuerySpecification(ctx *parsergen.QuerySpecificationContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitSelectSpec(ctx *parsergen.SelectSpecContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitSelectElements(ctx *parsergen.SelectElementsContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitSelectColumnElement(ctx *parsergen.SelectColumnElementContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitSelectFunctionElement(ctx *parsergen.SelectFunctionElementContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitFromClause(ctx *parsergen.FromClauseContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitLimitClause(ctx *parsergen.LimitClauseContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitLimitClauseAtom(ctx *parsergen.LimitClauseAtomContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitFullId(ctx *parsergen.FullIdContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitTableName(ctx *parsergen.TableNameContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitFullColumnName(ctx *parsergen.FullColumnNameContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitMysqlVariable(ctx *parsergen.MysqlVariableContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitCollationName(ctx *parsergen.CollationNameContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitUid(ctx *parsergen.UidContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitSimpleId(ctx *parsergen.SimpleIdContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitDottedId(ctx *parsergen.DottedIdContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitDecimalLiteral(ctx *parsergen.DecimalLiteralContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitStringLiteral(ctx *parsergen.StringLiteralContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitBooleanLiteral(ctx *parsergen.BooleanLiteralContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitHexadecimalLiteral(ctx *parsergen.HexadecimalLiteralContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitConstant(ctx *parsergen.ConstantContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitExpressions(ctx *parsergen.ExpressionsContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitAggregateFunctionCall(ctx *parsergen.AggregateFunctionCallContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitAggregateWindowedFunction(ctx *parsergen.AggregateWindowedFunctionContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitIsExpression(ctx *parsergen.IsExpressionContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitNotExpression(ctx *parsergen.NotExpressionContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitLogicalExpression(ctx *parsergen.LogicalExpressionContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitPredicateExpression(ctx *parsergen.PredicateExpressionContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitExpressionAtomPredicate(ctx *parsergen.ExpressionAtomPredicateContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitBinaryComparisonPredicate(ctx *parsergen.BinaryComparisonPredicateContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitInPredicate(ctx *parsergen.InPredicateContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitConstantExpressionAtom(ctx *parsergen.ConstantExpressionAtomContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitFullColumnNameExpressionAtom(ctx *parsergen.FullColumnNameExpressionAtomContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitFunctionCallExpressionAtom(ctx *parsergen.FunctionCallExpressionAtomContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitUnaryExpressionAtom(ctx *parsergen.UnaryExpressionAtomContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitBinaryExpressionAtom(ctx *parsergen.BinaryExpressionAtomContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitNestedExpressionAtom(ctx *parsergen.NestedExpressionAtomContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitUnaryOperator(ctx *parsergen.UnaryOperatorContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitComparisonOperator(ctx *parsergen.ComparisonOperatorContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *AstBuilder) VisitLogicalOperator(ctx *parsergen.LogicalOperatorContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func NewAstBuilder(opts ...AstBuilderOption) *AstBuilder {
	v := &AstBuilder{}
	v.apply(opts...)
	return v
}
