package goar

type EnumRelationalOperators int

const (
	_   EnumRelationalOperators = iota
	EQ                          // equal
	NE                          // not equal
	LT                          // less than
	LTE                         // less than or equal
	GT                          // greater than
	GTE                         // greater than or equal
	IN
)

type EnumLogicalOperators int

const (
	_ EnumLogicalOperators = iota
	OR
	AND
	NOT
)

type EnumSortOrders int

const (
	_ EnumSortOrders = iota
	ASC
	DESC
)

type EnumAggregations int

const (
	_ EnumAggregations = iota
	SUM
	GROUP
)

type Querier interface {
	Query() *Query
	SetQuery(*Query)
	Pluck(...interface{}) *ActiveRecord
	Where(QueryCondition) *ActiveRecord
	Order(OrderBy) *ActiveRecord
	Sum(fields ...interface{}) *ActiveRecord
	Distinct() *ActiveRecord
	//Or(QueryCondition) *ActiveRecord
	Run(results interface{}) error
}

type Query struct {
	//db              *DB
	//OrConditions    []QueryCondition
	//NotConditions   []QueryCondition
	Plucks          []interface{}
	WhereConditions []QueryCondition
	OrderBys        []OrderBy
	Joins           string
	Offset          string
	Limit           string
	Aggregations    map[EnumAggregations][]interface{}
	Distinct        bool
	err             error
}

type QueryCondition struct {
	LogicalOperator    EnumLogicalOperators
	Key                string
	RelationalOperator EnumRelationalOperators
	Value              interface{}
}

type OrderBy struct {
	Key       string
	SortOrder EnumSortOrders
}

func NewQuery() *Query {
	return &Query{
		Aggregations: make(map[EnumAggregations][]interface{}),
		Distinct:     false}
}

//func (q *Query) Where(where WhereCondition) *Query {
//q.WhereConditions = append(q.WhereConditions, where)
//return q
//}

//func (q *Query) Run(results interface{}) error {
//model := q.Model
//return model.(Persister).DbSearch(results)
//}
