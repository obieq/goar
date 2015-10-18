package mssql

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jinzhu/gorm"
	. "github.com/obieq/goar"
)

type ArMsSql struct {
	ActiveRecord
	ID int `gorm:"primary_key" json:"id,omitempty"`
	//ID string `sql:"type:varchar(36)" gorm:"primary_key" json:"id,omitempty"`
	Timestamps
}

// interface assertions
// https://splice.com/blog/golang-verify-type-implements-interface-compile-time/
var _ Persister = (*ArMsSql)(nil)
var _ RDBMSer = (*ArMsSql)(nil)

var (
	clients = map[string]gorm.DB{}
)

func connect(connName string, env string) (client gorm.DB) {
	c := Config
	if c == nil {
		log.Panic("goar config cannot be nil")
	}

	connKey := env + "_mssql_" + connName
	m, found := c.MSSQLDBs[connKey]
	if !found {
		log.Panic("mssql connection not found:", connKey)
	}

	connString := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s", m.Server, m.Port, m.DBName, m.Username, m.Password)

	if m.Debug {
		log.Println("connString:", connString)
	}

	db, err := gorm.Open("mssql", connString)
	if err != nil {
		log.Panic("Open mssql database failed:", err)
	}

	// set connection properties
	db.DB().SetMaxIdleConns(m.MaxIdleConnections)
	db.DB().SetMaxOpenConns(m.MaxOpenConnections)

	// set log mode
	db.LogMode(m.Debug)

	// test the connection
	if err = db.DB().Ping(); err != nil {
		log.Panic("mssql ping failed:", err.Error())
	}

	//return connection
	return db
}

func (ar *ArMsSql) SetKey(key string) {
	// TODO: set guid key here once that's implemented
	//ar.ID = key
}

func (ar *ArMsSql) Client() gorm.DB {
	self := ar.Self()
	connectionKey := self.DBConnectionName() + "_" + self.DBConnectionEnvironment()
	if self == nil {
		log.Panic("ar.Self() cannot be blank!")
	}

	conn, found := clients[connectionKey]
	if !found {
		conn = connect(self.DBConnectionName(), self.DBConnectionEnvironment())
		clients[connectionKey] = conn
	}

	return conn
}

func (ar *ArMsSql) All(models interface{}, opts map[string]interface{}) (err error) {
	var limit int = 100

	// set limit
	if opts["limit"] != nil {
		limit = opts["limit"].(int)
		if limit > 1000 { // max limit is 1000
			return errors.New("limit must be less than 1001")
		} else if limit < 1 {
			return errors.New("limit must be greater than 0")
		}
	}

	client := ar.Client()
	return client.Limit(limit).Find(models).Error
}

func (ar *ArMsSql) Truncate() (numRowsDeleted int, err error) {
	client := ar.Client()
	tblName := client.NewScope(ar.Self()).TableName()
	return -1, client.Exec("TRUNCATE TABLE " + tblName).Error
}

func (ar *ArMsSql) Find(id interface{}, out interface{}) error {
	client := ar.Client()
	return client.First(out, id).Error
}

func (ar *ArMsSql) DbSave() error {
	var err error

	//if ar.UpdatedAt != nil {
	//err = client.Save(ar.Self()).Error
	////_, err = client.Put(ar.ModelName(), ar.ID, ar.Self())
	//} else {
	//_, err = client.PutIfAbsent(ar.ModelName(), ar.ID, ar.Self())
	client := ar.Client()
	err = client.Create(ar.Self()).Error
	//}

	return err
}

func (ar *ArMsSql) DbDelete() (err error) {
	client := ar.Client()
	return client.Delete(ar.Self()).Error
}

func (ar *ArMsSql) DbSearch(models interface{}) (err error) {
	var query, sort string
	//var response *c.SearchResults
	//query := r.Db(DbName()).Table(ar.Self().ModelName())

	// plucks
	//query = processPlucks(query, ar)

	// where conditions
	if query, err = processWhereConditions(ar); err != nil {
		return err
	}

	// aggregations
	//if query, err = processAggregations(query, ar); err != nil {
	//return err
	//}

	// order bys
	sort = processSorts(ar)

	// TODO: delete!
	log.Printf("DbSearch query: %s", query)

	// run search
	if sort == "" {
		//if response, err = client.Search(ar.ModelName(), query, 100, 0); err != nil {
		//return err
		//}
	} else {
		//if response, err = client.SearchSorted(ar.ModelName(), query, sort, 100, 0); err != nil {
		//return err
		//}
	}

	//return mapResults(response.Results, models)
	return nil
}

func (ar *ArMsSql) SpExecResultSet(spName string, params map[string]interface{}, models interface{}) (err error) {
	client := ar.Client()
	if params == nil {
		return client.Raw("exec " + spName).Scan(models).Error
	} else {
		return client.Raw("exec " + spName + buildSpParams(params)).Scan(models).Error
	}
}

func buildSpParams(params map[string]interface{}) string {
	var kvs []string
	key := ""

	for k, v := range params {
		key = " @" + k + " = "

		switch v.(type) {
		case string:
			kvs = append(kvs, key+"'"+v.(string)+"'")
		case int:
			kvs = append(kvs, key+strconv.Itoa(v.(int)))
		case int64:
			kvs = append(kvs, key+strconv.FormatInt(v.(int64), 10))
		default:
			log.Panic("the following stored proc param type has not been implemented: ", reflect.TypeOf(v))
		}
	}

	return strings.Join(kvs, ",")
}

//func processPlucks(query r.Term, ar *ArRethinkDb) r.Term {
//if plucks := ar.Query().Plucks; plucks != nil {
//query = query.Pluck(plucks...)
//}

//return query
//}

func mapResults(orchestrateResults interface{}, models interface{}) (err error) {
	// now, map orchstrate's raw json to the desired active record type
	//modelsv := reflect.ValueOf(models)
	//if modelsv.Kind() != reflect.Ptr || modelsv.Elem().Kind() != reflect.Slice {
	//panic("models argument must be a slice address")
	//}
	//slicev := modelsv.Elem()
	//elemt := slicev.Type().Elem()

	//switch t := orchestrateResults.(type) {
	//case []c.KVResult:
	//for _, result := range t {
	//elemp := reflect.New(elemt)
	//if err = result.Value(elemp.Interface()); err != nil {
	//return err
	//}

	//slicev = reflect.Append(slicev, elemp.Elem())
	//}
	//case []c.SearchResult:
	//for _, result := range t {
	//elemp := reflect.New(elemt)
	//if err = result.Value(elemp.Interface()); err != nil {
	//return err
	//}

	//slicev = reflect.Append(slicev, elemp.Elem())
	//}
	//default:
	//return errors.New(fmt.Sprintf("Orchestrate Response Type Not Mapped: %v", t))
	//}

	//// assign mapped results to the caller's supplied array
	//modelsv.Elem().Set(slicev)

	//return err
	return nil
}

func processWhereConditions(ar *ArMsSql) (query string, err error) {
	var whereStmt, whereCondition string

	if len(ar.Query().WhereConditions) > 0 {
		for index, where := range ar.Query().WhereConditions {
			switch where.RelationalOperator {
			case EQ: // equal
				whereCondition = where.Key + ":" + fmt.Sprintf("%v", where.Value)
				//whereCondition = where.Key + ":" + where.Value.(string)
				//whereCondition = r.Row.Field(where.Key).Eq(where.Value)
			//case NE: // not equal
			//whereCondition = r.Row.Field(where.Key).Ne(where.Value)
			//case LT: // less than
			//whereCondition = r.Row.Field(where.Key).Lt(where.Value)
			//case LTE: // less than or equal
			//whereCondition = r.Row.Field(where.Key).Le(where.Value)
			//case GT: // greater than
			//// TODO: create function to set range based on type???
			//whereCondition = where.Key + ":[" + fmt.Sprintf("%v", where.Value) + " TO *]"
			//whereCondition = r.Row.Field(where.Key).Gt(where.Value)
			case GTE: // greater than or equal
				whereCondition = where.Key + ":[" + fmt.Sprintf("%v", where.Value) + " TO *]"
			//whereCondition = r.Row.Field(where.Key).Ge(where.Value)
			// case IN: // TODO: implement!!!!
			default:
				return query, errors.New(fmt.Sprintf("invalid comparison operator: %v", where.RelationalOperator))
			}

			if index == 0 {
				whereStmt = whereCondition
				//if where.LogicalOperator == NOT {
				//whereStmt = whereStmt.Not()
				//}
			} else {
				switch where.LogicalOperator {
				case AND:
					whereStmt = whereStmt + " AND " + whereCondition
					//whereStmt = whereStmt.And(whereCondition)
				case OR:
					whereStmt = whereStmt + " OR " + whereCondition
				//whereStmt = whereStmt.Or(whereCondition)
				////case NOT:
				////whereStmt = whereStmt.And(whereCondition).Not()
				default:
					whereStmt = whereStmt + " AND " + whereCondition
					//whereStmt = whereStmt.And(whereCondition)
				}
			}
		}

		// TODO: delete!!
		log.Printf("DbSearch whereStmt: %s", whereStmt)
		//query = query.Filter(whereStmt)
		//query = query.Filter(whereStmt)
	}

	return whereStmt, nil
}

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

func processSorts(ar *ArMsSql) (sort string) {
	if len(ar.Query().OrderBys) > 0 {
		sort = ""

		for i, orderBy := range ar.Query().OrderBys {
			if i > 0 {
				sort += ","
			}

			sort += "value." + orderBy.Key + ":"

			switch orderBy.SortOrder {
			case DESC: // descending
				sort += "desc"
			default: // ascending
				sort += "asc"
			}
		}
	}

	return sort
}
