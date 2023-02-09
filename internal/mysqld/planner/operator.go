package planner

type UnaryOperator = int

const (
	UnaryOperator_Unknown           UnaryOperator = iota
	UnaryOperator_ExclamationSymbol               // '!'
	UnaryOperator_Tilde                           // '~'
	UnaryOperator_Positive                        // '+'
	UnaryOperator_Negative                        // '-'
	UnaryOperator_Not                             // 'not'
)

type ComparisonOperator = int

const (
	ComparisonOperator_Unknown      ComparisonOperator = 0
	ComparisonOperator_Equal        ComparisonOperator = 1                                                         // '='
	ComparisonOperator_GreaterThan  ComparisonOperator = 2                                                         // '>'
	ComparisonOperator_LessThan     ComparisonOperator = 4                                                         // '<'
	ComparisonOperator_NotEqual     ComparisonOperator = 8                                                         // '<>', '!='
	ComparisonOperator_LessEqual    ComparisonOperator = ComparisonOperator_LessThan | ComparisonOperator_Equal    // '<='
	ComparisonOperator_GreaterEqual ComparisonOperator = ComparisonOperator_GreaterThan | ComparisonOperator_Equal // '>='
)

type LogicalOperator = int

const (
	LogicalOperator_Unknown LogicalOperator = iota
	LogicalOperator_And
	LogicalOperator_Or
)

type IsOperator = int

const (
	IsOperator_Unknown IsOperator = iota
	IsOperator_Is
	IsOperator_IsNot
)

type InOperator = int

const (
	InOperator_Unknown InOperator = iota
	InOperator_In
	InOperator_NotIn
)
