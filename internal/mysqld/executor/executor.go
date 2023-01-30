package executor

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/milvus-io/milvus/internal/util/typeutil"

	"github.com/milvus-io/milvus-proto/go-api/schemapb"

	"github.com/milvus-io/milvus-proto/go-api/commonpb"
	"github.com/milvus-io/milvus-proto/go-api/milvuspb"
	"github.com/milvus-io/milvus/internal/util/commonpbutil"
	querypb "github.com/xelabs/go-mysqlstack/sqlparser/depends/query"

	"github.com/milvus-io/milvus/internal/mysqld/planner"
	"github.com/milvus-io/milvus/internal/types"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/xelabs/go-mysqlstack/sqlparser/depends/sqltypes"
)

type Executor interface {
	Run(plan *planner.PhysicalPlan) (*sqltypes.Result, error)
}

// defaultExecutor only translates sql to rpc. TODO: Better to use vacalno model or batch model.
type defaultExecutor struct {
	s types.ProxyComponent
}

func (e *defaultExecutor) Run(plan *planner.PhysicalPlan) (*sqltypes.Result, error) {
	return e.dispatch(plan)
}

func (e *defaultExecutor) dispatch(plan *planner.PhysicalPlan) (*sqltypes.Result, error) {
	switch stmt := plan.Stmt.(type) {
	case *ast.SelectStmt:
		return e.runSelectStmt(stmt)
	default:
		return nil, fmt.Errorf("not supported sql: %s", stmt.OriginalText())
	}
}

func (e *defaultExecutor) runSelectStmt(stmt *ast.SelectStmt) (*sqltypes.Result, error) {
	// select 4 + 5;
	if stmt.From == nil {
		return e.runStmtWithoutFrom(stmt)
	}

	// select count(*) from t;
	if stmt.Where == nil {
		return e.runStmtWithoutFilter(stmt)
	}

	return e.runStmtWithFilter(stmt)
}

func (e *defaultExecutor) runStmtWithoutFrom(stmt *ast.SelectStmt) (*sqltypes.Result, error) {
	return nil, fmt.Errorf("sql didn't specify the table name: %s", stmt.OriginalText())
}

func (e *defaultExecutor) runStmtWithoutFilter(stmt *ast.SelectStmt) (*sqltypes.Result, error) {
	if err := matchCountRule(stmt); err != nil {
		return nil, err
	}

	tableName, err := getTableName(stmt)
	if err != nil {
		return nil, err
	}

	rowCnt, err := e.getCount(tableName)
	if err != nil {
		return nil, err
	}

	result1 := wrapCountResult(rowCnt, "count(*)")
	return result1, nil
}

func matchCountRule(stmt *ast.SelectStmt) error {
	if len(stmt.Fields.Fields) != 1 {
		return fmt.Errorf("only one count expr is allowed without filter: %s", stmt.OriginalText())
	}

	field := stmt.Fields.Fields[0]

	// select * from collection;
	if field.WildCard != nil {
		return fmt.Errorf("query without filter is not allowed: %s", stmt.OriginalText())
	}

	aggFuncExpr, ok := field.Expr.(*ast.AggregateFuncExpr)
	if !ok || aggFuncExpr.F != ast.AggFuncCount {
		return fmt.Errorf("only count aggregation is allowed: %s", stmt.OriginalText())
	}

	// handle arguments. For example, count(*), count(column). TODO: handle `distinct` argument.
	args := aggFuncExpr.Args
	if len(args) != 1 {
		return fmt.Errorf("count complex expression is not allowed: %s", stmt.OriginalText())
	}

	return nil
}

func (e *defaultExecutor) runStmtWithFilter(stmt *ast.SelectStmt) (*sqltypes.Result, error) {
	res, err := e.tryPointGetOnly(stmt)
	if err == nil {
		return res, nil
	}
	return e.tryCountWithFilter(stmt)
}

func (e *defaultExecutor) tryPointGetOnly(stmt *ast.SelectStmt) (*sqltypes.Result, error) {
	tableName, err := getTableName(stmt)
	if err != nil {
		return nil, err
	}

	filter, err := newTextRestorer().getText(stmt.Where)
	if err != nil {
		return nil, err
	}

	outputs, err := getOutputFields(stmt.Fields)
	if err != nil {
		return nil, err
	}

	res, err := e.query(tableName, filter, outputs)
	if err != nil {
		return nil, err
	}

	return wrapQueryResults(res), nil
}

func (e *defaultExecutor) tryCountWithFilter(stmt *ast.SelectStmt) (*sqltypes.Result, error) {
	if err := matchCountRule(stmt); err != nil {
		return nil, err
	}

	tableName, err := getTableName(stmt)
	if err != nil {
		return nil, err
	}

	filter, err := newTextRestorer().getText(stmt.Where)
	if err != nil {
		return nil, err
	}

	// TODO: check if `*` match vector field.
	outputs := []string{"*"}

	res, err := e.query(tableName, filter, outputs)
	if err != nil {
		return nil, err
	}

	nColumn := len(res.GetFieldsData())
	nRow := 0
	if nColumn > 0 {
		nRow = typeutil.GetRowCount(res.GetFieldsData()[0])
	}

	return wrapCountResult(nRow, "count(*)"), nil
}

func getTableName(stmt *ast.SelectStmt) (string, error) {
	from := stmt.From

	if from.TableRefs.Tp != 0 || from.TableRefs.Right != nil {
		return "", fmt.Errorf("invalid data source: %s", stmt.OriginalText())
	}

	if from.TableRefs.Left == nil {
		return "", fmt.Errorf("invalid data source: %s", stmt.OriginalText())
	}

	left := from.TableRefs.Left
	tableSourceNode, ok := left.(*ast.TableSource)
	if !ok {
		return "", fmt.Errorf("invalid data source: %s", stmt.OriginalText())
	}
	tableNameNode, ok := tableSourceNode.Source.(*ast.TableName)
	if !ok {
		return "", fmt.Errorf("invalid data source: %s", stmt.OriginalText())
	}

	return tableNameNode.Name.String(), nil
}

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

func (e *defaultExecutor) getCount(tableName string) (int, error) {
	req := &milvuspb.GetCollectionStatisticsRequest{
		Base:           commonpbutil.NewMsgBase(),
		CollectionName: tableName,
	}
	resp, err := e.s.GetCollectionStatistics(context.TODO(), req)
	if err != nil {
		return 0, err
	}
	if resp.GetStatus().GetErrorCode() != commonpb.ErrorCode_Success {
		return 0, errors.New(resp.GetStatus().GetReason())
	}
	rowCnt, err := strconv.Atoi(resp.GetStats()[0].GetValue())
	if err != nil {
		return 0, err
	}
	return rowCnt, nil
}

func wrapCountResult(rowCnt int, column string) *sqltypes.Result {
	result1 := &sqltypes.Result{
		Fields: []*querypb.Field{
			{
				Name: column,
				Type: querypb.Type_INT64,
			},
		},
		Rows: [][]sqltypes.Value{
			{
				sqltypes.NewInt64(int64(rowCnt)),
			},
		},
	}
	return result1
}

func (e *defaultExecutor) query(tableName string, filter string, outputs []string) (*milvuspb.QueryResults, error) {
	req := &milvuspb.QueryRequest{
		Base:               commonpbutil.NewMsgBase(),
		DbName:             "",
		CollectionName:     tableName,
		Expr:               filter,
		OutputFields:       outputs,
		PartitionNames:     nil,
		TravelTimestamp:    0,
		GuaranteeTimestamp: 0,
		QueryParams:        nil,
	}

	resp, err := e.s.Query(context.TODO(), req)
	if err != nil {
		return nil, err
	}

	if resp.GetStatus().GetErrorCode() != commonpb.ErrorCode_Success {
		return nil, errors.New(resp.GetStatus().GetReason())
	}

	return resp, nil
}

func wrapQueryResults(res *milvuspb.QueryResults) *sqltypes.Result {
	fieldsData := res.GetFieldsData()
	nColumn := len(fieldsData)
	fields := make([]*querypb.Field, 0, nColumn)

	if nColumn <= 0 {
		return &sqltypes.Result{}
	}

	for i := 0; i < nColumn; i++ {
		fields = append(fields, getSQLField(res.GetCollectionName(), fieldsData[i]))
	}

	nRow := typeutil.GetRowCount(fieldsData[0])
	rows := make([][]sqltypes.Value, 0, nRow)
	for i := 0; i < nRow; i++ {
		row := make([]sqltypes.Value, 0, nColumn)
		for j := 0; j < nColumn; j++ {
			row = append(row, getDataSingle(fieldsData[j], i))
		}
		rows = append(rows, row)
	}

	return &sqltypes.Result{
		Fields: fields,
		Rows:   rows,
	}
}

func getSQLField(tableName string, fieldData *schemapb.FieldData) *querypb.Field {
	return &querypb.Field{
		Name:         fieldData.GetFieldName(),
		Type:         toSQLType(fieldData.GetType()),
		Table:        tableName,
		OrgTable:     "",
		Database:     "",
		OrgName:      "",
		ColumnLength: 0,
		Charset:      0,
		Decimals:     0,
		Flags:        0,
	}
}

func toSQLType(t schemapb.DataType) querypb.Type {
	switch t {
	case schemapb.DataType_Bool:
		// TODO: tinyint
		return querypb.Type_UINT8
	case schemapb.DataType_Int8:
		return querypb.Type_INT8
	case schemapb.DataType_Int16:
		return querypb.Type_INT16
	case schemapb.DataType_Int32:
		return querypb.Type_INT32
	case schemapb.DataType_Int64:
		return querypb.Type_INT64
	case schemapb.DataType_Float:
		return querypb.Type_FLOAT32
	case schemapb.DataType_Double:
		return querypb.Type_FLOAT64
	case schemapb.DataType_VarChar:
		return querypb.Type_VARCHAR
		// TODO: vector.
	default:
		return querypb.Type_NULL_TYPE
	}
}

func getDataSingle(fieldData *schemapb.FieldData, idx int) sqltypes.Value {
	switch fieldData.GetType() {
	case schemapb.DataType_Bool:
		// TODO: tinyint
		return sqltypes.NewInt32(1)
	case schemapb.DataType_Int8:
		v := fieldData.Field.(*schemapb.FieldData_Scalars).Scalars.Data.(*schemapb.ScalarField_IntData).IntData.GetData()[idx]
		return sqltypes.MakeTrusted(sqltypes.Int8, strconv.AppendInt(nil, int64(v), 10))
	case schemapb.DataType_Int16:
		v := fieldData.Field.(*schemapb.FieldData_Scalars).Scalars.Data.(*schemapb.ScalarField_IntData).IntData.GetData()[idx]
		return sqltypes.MakeTrusted(sqltypes.Int16, strconv.AppendInt(nil, int64(v), 10))
	case schemapb.DataType_Int32:
		v := fieldData.Field.(*schemapb.FieldData_Scalars).Scalars.Data.(*schemapb.ScalarField_IntData).IntData.GetData()[idx]
		return sqltypes.MakeTrusted(sqltypes.Int32, strconv.AppendInt(nil, int64(v), 10))
	case schemapb.DataType_Int64:
		v := fieldData.Field.(*schemapb.FieldData_Scalars).Scalars.Data.(*schemapb.ScalarField_LongData).LongData.GetData()[idx]
		return sqltypes.MakeTrusted(sqltypes.Int64, strconv.AppendInt(nil, v, 10))
	case schemapb.DataType_Float:
		v := fieldData.Field.(*schemapb.FieldData_Scalars).Scalars.Data.(*schemapb.ScalarField_FloatData).FloatData.GetData()[idx]
		return sqltypes.MakeTrusted(sqltypes.Float32, strconv.AppendFloat(nil, float64(v), 'f', -1, 64))
	case schemapb.DataType_Double:
		v := fieldData.Field.(*schemapb.FieldData_Scalars).Scalars.Data.(*schemapb.ScalarField_DoubleData).DoubleData.GetData()[idx]
		return sqltypes.MakeTrusted(sqltypes.Float64, strconv.AppendFloat(nil, v, 'g', -1, 64))
	case schemapb.DataType_VarChar:
		v := fieldData.Field.(*schemapb.FieldData_Scalars).Scalars.Data.(*schemapb.ScalarField_StringData).StringData.GetData()[idx]
		return sqltypes.NewVarChar(v)

		// TODO: vector.
	default:
		// TODO: should raise error here.
		return sqltypes.NewInt32(1)
	}
}

func NewDefaultExecutor(s types.ProxyComponent) Executor {
	return &defaultExecutor{s: s}
}
