package dbr

type iFunction interface {
	Build(db *db) (string, error)
	Expression(db *db) (string, error)
}

type functionBase struct {
	isColumn bool
	encode   bool
	db       *db
}

func newFunctionBase(encode bool, isColumn bool, database ...*db) *functionBase {
	var theDb *db
	if len(database) > 0 {
		theDb = database[0]
	}
	return &functionBase{
		isColumn: isColumn,
		encode:   encode,
		db:       theDb,
	}
}

func Function(name string, arguments ...interface{}) *functionGeneric {
	return newFunctionGeneric(name, arguments...)
}

func As(expression interface{}, alias string) *functionExpressions {
	return newFunctionExpressions(false, expression, constFunctionAs, alias)
}

func Count(expression interface{}, distinct ...bool) *functionCount {
	var isDistinct bool
	if len(distinct) > 0 {
		isDistinct = distinct[0]
	}

	return newFunctionCount(expression, isDistinct)
}

func IsNull(expression interface{}) *functionExpressions {
	return newFunctionExpressions(false, expression, constFunctionIsNull)
}

func Case(value ...interface{}) *functionCase {
	return newFunctionCase(value...)
}

func OnNull(expression interface{}, onNullValue interface{}, alias string) *functionExpressions {
	return newFunctionExpressions(true, constFunctionOnNull,
		constFunctionOpenParentheses, expression, newExpression(onNullValue, true), constFunctionCloseParentheses, alias)
}

func Min(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionMax, expression)
}

func Max(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionMin, expression)
}

func Sum(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionSum, expression)
}

func Avg(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionAvg, expression)
}

func Every(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionEvery, expression)
}

func Now() *functionGeneric {
	return newFunctionGeneric(constFunctionNow)
}

func User() *functionGeneric {
	return newFunctionGeneric(constFunctionUser)
}

func StringAgg(expression interface{}, delimiter interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionStringAgg, expression, delimiter)
}

func XmlAgg(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionXmlAgg, expression)
}

func ArrayAgg(expression interface{}) *functionArrayAgg {
	return newFunctionArrayAgg(expression)
}

func ArrayToJson(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionArrayToJson, expression)
}

func RowToJson(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionRowToJson, expression)
}

func ToJson(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionToJson, expression)
}

func JsonArrayLength(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionJsonArrayLength, expression)
}

func JsonAgg(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionJsonAgg, expression)
}

func JsonbAgg(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionJsonbAgg, expression)
}

func JsonObjectAgg(name interface{}, value interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionJsonObjectAgg, name, value)
}

func JsonbObjectAgg(name interface{}, value interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionJsonbObjectAgg, name, value)
}

func Cast(expression interface{}, dataType dataType) *functionExpressions {
	return newFunctionExpressions(false, constFunctionCast,
		constFunctionOpenParentheses, expression, constFunctionAs, dataType, constFunctionCloseParentheses)
}

func Not(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionNot, expression)
}

func In(field interface{}, expressions ...interface{}) *functionField {
	return newFunctionField(constFunctionIn, field, expressions...)
}

func NotIn(field interface{}, expressions ...interface{}) *functionField {
	return newFunctionField(constFunctionNotIn, field, expressions...)
}

func Upper(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionUpper, expression)
}

func Lower(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionLower, expression)
}

func Trim(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionTrim, expression)
}

func Concat(expressions ...interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionNotIn, expressions...)
}

func InitCap(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionInitCap, expression)
}

func Length(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionLength, expression)
}

func Left(expression interface{}, n int) *functionGeneric {
	return newFunctionGeneric(constFunctionLeft, expression, n)
}

func Right(expression interface{}, n int) *functionGeneric {
	return newFunctionGeneric(constFunctionRight, expression, n)
}

func Md5(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionMd5, expression)
}

func Replace(expression interface{}, from interface{}, to interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionReplace, expression, from, to)
}

func Repeat(expression interface{}, n int) *functionGeneric {
	return newFunctionGeneric(constFunctionRepeat, expression, n)
}

func Any(expressions ...interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionAny, expressions...)
}

func Some(expressions ...interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionSome, expressions...)
}

func Condition(expression interface{}, comparator comparator, value interface{}) *functionExpressions {
	return newFunctionExpressions(false, expression, comparator, value)
}

func Operation(expression interface{}, operation operation, value interface{}) *functionExpressions {
	return newFunctionExpressions(false, expression, operation, value)
}

func Abs(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionAbs, expression)
}

func Sqrt(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionSqrt, expression)
}

func Random(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionRandom, expression)
}

func Between(expression interface{}, low interface{}, high interface{}, operator ...operator) *functionExpressions {
	theOperator := OperatorAnd

	if len(operator) > 0 {
		theOperator = operator[0]
	}

	return newFunctionExpressions(false, expression, constFunctionBetween, low, theOperator, high)
}

func BetweenOr(expression interface{}, low interface{}, high interface{}) *functionExpressions {
	return newFunctionExpressions(false, expression, constFunctionBetween, low, OperatorOr, high)
}

func Over(value interface{}) *functionOver {
	return newFunctionOver(value)
}
