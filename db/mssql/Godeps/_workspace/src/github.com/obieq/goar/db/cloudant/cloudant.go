package cloudant

import (
	"os"
	"reflect"
	"time"

	c "github.com/obieq/go-cloudant"
	. "github.com/obieq/goar"
)

type ArCloudant struct {
	ActiveRecord
	c.CloudantDocument
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	response  c.CloudantDocumentResponse
}

var (
	client   *c.Client
	database *c.Database
	dbName   string
)

var connectOpts = func() map[string]string {
	opts := make(map[string]string)
	opts["base_uri"] = os.Getenv("CLOUDANT_URL")
	opts["api_key"] = os.Getenv("CLOUDANT_API_KEY")
	opts["api_password"] = os.Getenv("CLOUDANT_API_PASSWORD")

	return opts
}

func connect() (client *c.Client, err error) {
	opts := connectOpts()
	err = nil

	return c.NewClient(opts["base_uri"], opts["api_key"], opts["api_password"]), nil
}

func init() {
	client, _ = connect()
}

func Client() *c.Client {
	return client
}

func Database() *c.Database {
	return database
}

func DbName() string {
	return dbName
}

func SetDbName(name string) {
	dbName = name
	database = c.NewDatabase(name, client)
}

func (ar *ArCloudant) All(results interface{}) error {
	//rows, err := r.Db(dbName).Table(ar.Self().ModelName()).Run(session)
	//if err != nil {
	//log.Println(err)
	//} else {
	//err = rows.All(results)
	////modelInterface := reflect.New(modelVal.Type()).Interface()
	////for rows.Next(&modelInterface) {
	////result = append(result, modelInterface)
	////}
	////for rows.Next() {
	////modelInterface := reflect.New(modelVal.Type()).Interface()
	////err = rows.Scan(&modelInterface)
	////if err == nil { // would like to break if err 1= nil, but then difficult to get 100% test coverage
	////result = append(result, modelInterface)
	////}
	////}
	//}

	////return result, err
	//return err

	return nil
}

//var truncate = func(modelName string) (*r.Cursor, error) {
//return r.Db(dbName).Table(modelName).Delete().Run(session)
//}

func (ar *ArCloudant) Truncate() (numRowsDeleted int, err error) {
	//if _, err = truncate(ar.Self().ModelName()); err != nil {
	//log.Println(err)
	//}

	//return 0, err

	return 0, nil
}

func (ar *ArCloudant) Find(id interface{}) (interface{}, error) {
	self := ar.Self()
	modelVal := reflect.ValueOf(self).Elem()
	modelInterface := reflect.New(modelVal.Type()).Interface()

	err := database.GetDocument(id.(string), &modelInterface)

	return modelInterface, err
}

func (ar *ArCloudant) DbSave() error {
	var err error
	var cdr c.CloudantDocumentResponse

	if ar.UpdatedAt != nil {
		cdr, err = database.UpdateDocument(ar, false)
	} else {
		cdr, err = database.CreateDocument(ar, false)
	}

	ar.response = cdr

	return err
}

func (ar *ArCloudant) DbDelete() (err error) {
	//self := ar.Self()
	//modelVal := reflect.ValueOf(self).Elem()
	//_, err = r.Db(dbName).Table(self.ModelName()).Get(modelVal.FieldByName("Id").Interface()).Delete().Run(session) // TODO: use PrimaryKey
	//return err
	return nil
}

func (ar *ArCloudant) DbSearch(results interface{}) (err error) {
	//query := r.Db(DbName()).Table(ar.Self().ModelName())

	//// plucks
	//query = processPlucks(query, ar)

	//// where conditions
	//if query, err = processWhereConditions(query, ar); err != nil {
	//return err
	//}

	//// aggregations
	//if query, err = processAggregations(query, ar); err != nil {
	//return err
	//}

	//// order bys
	//query = processOrderBys(query, ar)

	//// TODO: delete!
	//log.Printf("DbSearch query: %s", query)

	//rows, err := query.Run(Session())
	//if err != nil {
	//return err
	//}

	//return rows.All(results)
	return nil
}

//func processPlucks(query r.Term, ar *ArRethinkDb) r.Term {
//if plucks := ar.Query().Plucks; plucks != nil {
//query = query.Pluck(plucks...)
//}

//return query
//}

//func processWhereConditions(query r.Term, ar *ArRethinkDb) (r.Term, error) {
//var whereStmt, whereCondition r.Term

//if len(ar.Query().WhereConditions) > 0 {
//for index, where := range ar.Query().WhereConditions {
//switch where.RelationalOperator {
//case EQ: // equal
//whereCondition = r.Row.Field(where.Key).Eq(where.Value)
//case NE: // not equal
//whereCondition = r.Row.Field(where.Key).Ne(where.Value)
//case LT: // less than
//whereCondition = r.Row.Field(where.Key).Lt(where.Value)
//case LTE: // less than or equal
//whereCondition = r.Row.Field(where.Key).Le(where.Value)
//case GT: // greater than
//whereCondition = r.Row.Field(where.Key).Gt(where.Value)
//case GTE: // greater than or equal
//whereCondition = r.Row.Field(where.Key).Ge(where.Value)
//default:
//return query, errors.New(fmt.Sprintf("invalid comparison operator: %v", where.RelationalOperator))
//}

//if index == 0 {
//whereStmt = whereCondition
////if where.LogicalOperator == NOT {
////whereStmt = whereStmt.Not()
////}
//} else {
//switch where.LogicalOperator {
//case AND:
//whereStmt = whereStmt.And(whereCondition)
//case OR:
//whereStmt = whereStmt.Or(whereCondition)
////case NOT:
////whereStmt = whereStmt.And(whereCondition).Not()
//default:
//whereStmt = whereStmt.And(whereCondition)
//}
//}
//}

//// TODO: delete!!
//log.Printf("DbSearch whereStmt: %s", whereStmt)
//query = query.Filter(whereStmt)
//}

//return query, nil
//}

//func processAggregations(query r.Term, ar *ArRethinkDb) (r.Term, error) {
//// sum
//if sum := ar.Query().Aggregations[SUM]; sum != nil {
//if len(sum) == 1 {
//query = query.Sum(sum...)
//} else {
//return query, errors.New(fmt.Sprintf("rethinkdb does not support summing more than one field at a time: %v", sum))
//}
//}

//// distinct
//if ar.Query().Distinct {
//query = query.Distinct()
//}

//return query, nil
//}

//func processOrderBys(query r.Term, ar *ArRethinkDb) r.Term {
//if len(ar.Query().OrderBys) > 0 {
//orderBys := []interface{}{}

//for _, orderBy := range ar.Query().OrderBys {
//switch orderBy.SortOrder {
//case DESC: // descending
//orderBys = append(orderBys, r.Desc(orderBy.Key))
//default: // ascending
//orderBys = append(orderBys, r.Asc(orderBy.Key))
//}
//}

//query = query.OrderBy(orderBys...)
//}

//return query
//}
