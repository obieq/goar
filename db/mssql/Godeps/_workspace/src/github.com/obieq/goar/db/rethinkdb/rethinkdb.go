package rethinkdb

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"

	r "github.com/dancannon/gorethink"
	goar "github.com/obieq/goar"
)

type ArRethinkDb struct {
	goar.ActiveRecord
	ID        string    `gorethink:"id,omitempty"`
	CreatedAt time.Time `gorethink:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorethink:"updated_at,omitempty" json:"updated_at,omitempty"`
}

// interface assertions
// https://splice.com/blog/golang-verify-type-implements-interface-compile-time/
var _ goar.Persister = (*ArRethinkDb)(nil)

var (
	clients = map[string]*r.Session{}
)

// this facilitates integration/unit testing
var connOpts = func(opts *goar.RethinkDBConfig) r.ConnectOpts {
	return r.ConnectOpts{
		Addresses:     opts.Addresses,
		Database:      opts.DBName,
		MaxIdle:       opts.MaxIdleConnections,
		MaxOpen:       opts.MaxOpenConnections,
		AuthKey:       opts.AuthKey,
		DiscoverHosts: opts.DiscoverHosts,
	}
}

func connect(connName string, env string) (s *r.Session) {
	var err error
	c := goar.Config
	if c == nil {
		log.Panicln("goar config cannot be nil")
	}

	connKey := env + "_rethinkdb_" + connName
	m, found := c.RethinkDBs[connKey]
	if !found {
		log.Panicln("rethinkdb connection not found:", connKey)
	}

	if s, err = r.Connect(connOpts(m)); err != nil {
		log.Println("an error occurred while trying to connection to rethinkdb.  will now panic.", err)
		log.Panicln("could not connect to rethinkdb:", err)
	}

	return s
}

func (ar *ArRethinkDb) Client() *r.Session {
	self := ar.Self()
	if self == nil {
		log.Panicln("ar.Self() cannot be blank!")
	}

	connectionKey := self.DBConnectionName() + "_" + self.DBConnectionEnvironment()
	conn, found := clients[connectionKey]
	if !found {
		conn = connect(self.DBConnectionName(), self.DBConnectionEnvironment())
		clients[connectionKey] = conn
	}

	return conn
}

func (ar *ArRethinkDb) SetKey(key string) {
	ar.ID = key
}

func (ar *ArRethinkDb) All(results interface{}, opts map[string]interface{}) error {
	//result := []interface{}{}
	//self := ar.Self()
	//modelVal := reflect.ValueOf(self).Elem()
	rows, err := r.Table(ar.Self().ModelName()).Run(ar.Client())
	if err != nil {
		log.Println(err)
	} else {
		err = rows.All(results)
		//modelInterface := reflect.New(modelVal.Type()).Interface()
		//for rows.Next(&modelInterface) {
		//result = append(result, modelInterface)
		//}
		//for rows.Next() {
		//modelInterface := reflect.New(modelVal.Type()).Interface()
		//err = rows.Scan(&modelInterface)
		//if err == nil { // would like to break if err 1= nil, but then difficult to get 100% test coverage
		//result = append(result, modelInterface)
		//}
		//}
	}

	//return result, err
	return err
}

var truncate = func(session *r.Session, modelName string) (*r.Cursor, error) {
	return r.Table(modelName).Delete().Run(session)
}

func (ar *ArRethinkDb) Truncate() (numRowsDeleted int, err error) {
	if _, err = truncate(ar.Client(), ar.Self().ModelName()); err != nil {
		log.Println(err)
	}

	return 0, err
}

func (ar *ArRethinkDb) Find(id interface{}, out interface{}) error {
	row, err := r.Table(ar.ModelName()).Get(id).Run(ar.Client())

	if err != nil {
		log.Println(err)
	} else {
		if row.IsNil() { // return a not found error
			err = errors.New("record not found")
			log.Println("record not found for key:", id)
		} else {
			err = row.One(out)
		}
	}

	return err
}

func (ar *ArRethinkDb) DbSave() error {
	// Conflict parameter values: "error" (default), "replace", "update"
	// http://rethinkdb.com/api/javascript/insert/
	rslt, err := r.Table(ar.Self().ModelName()).Insert(ar.Self(), r.InsertOpts{Conflict: "update"}).RunWrite(ar.Client())
	if err == nil && ar.ID == "" { // if the client doesn't specify the PK, then Rethink will auto-generate it
		ar.ID = rslt.GeneratedKeys[0]
	}

	return err
}

func (ar *ArRethinkDb) DbDelete() (err error) {
	self := ar.Self()
	modelVal := reflect.ValueOf(self).Elem()
	_, err = r.Table(self.ModelName()).Get(modelVal.FieldByName("ID").Interface()).Delete().Run(ar.Client()) // TODO: use PrimaryKey

	return err
}

func (ar *ArRethinkDb) DbSearch(results interface{}) (err error) {
	query := r.Table(ar.Self().ModelName())

	// plucks
	query = processPlucks(query, ar)

	// where conditions
	if query, err = processWhereConditions(query, ar); err != nil {
		return err
	}

	// aggregations
	if query, err = processAggregations(query, ar); err != nil {
		return err
	}

	// order bys
	query = processOrderBys(query, ar)

	// TODO: delete!
	log.Printf("DbSearch query: %s", query)

	rows, err := query.Run(ar.Client())
	if err != nil {
		return err
	}

	return rows.All(results)
}

func processPlucks(query r.Term, ar *ArRethinkDb) r.Term {
	if plucks := ar.Query().Plucks; plucks != nil {
		query = query.Pluck(plucks...)
	}

	return query
}

func processWhereConditions(query r.Term, ar *ArRethinkDb) (r.Term, error) {
	var whereStmt, whereCondition r.Term

	if len(ar.Query().WhereConditions) > 0 {
		for index, where := range ar.Query().WhereConditions {
			switch where.RelationalOperator {
			case goar.EQ: // equal
				whereCondition = r.Row.Field(where.Key).Eq(where.Value)
			case goar.NE: // not equal
				whereCondition = r.Row.Field(where.Key).Ne(where.Value)
			case goar.LT: // less than
				whereCondition = r.Row.Field(where.Key).Lt(where.Value)
			case goar.LTE: // less than or equal
				whereCondition = r.Row.Field(where.Key).Le(where.Value)
			case goar.GT: // greater than
				whereCondition = r.Row.Field(where.Key).Gt(where.Value)
			case goar.GTE: // greater than or equal
				whereCondition = r.Row.Field(where.Key).Ge(where.Value)
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
				case goar.AND:
					whereStmt = whereStmt.And(whereCondition)
				case goar.OR:
					whereStmt = whereStmt.Or(whereCondition)
				//case goar.NOT:
				//whereStmt = whereStmt.And(whereCondition).Not()
				default:
					whereStmt = whereStmt.And(whereCondition)
				}
			}
		}

		// TODO: delete!!
		log.Printf("DbSearch whereStmt: %s", whereStmt)
		query = query.Filter(whereStmt)
	}

	return query, nil
}

func processAggregations(query r.Term, ar *ArRethinkDb) (r.Term, error) {
	// sum
	if sum := ar.Query().Aggregations[goar.SUM]; sum != nil {
		if len(sum) == 1 {
			query = query.Sum(sum...)
		} else {
			return query, errors.New(fmt.Sprintf("rethinkdb does not support summing more than one field at a time: %v", sum))
		}
	}

	// distinct
	if ar.Query().Distinct {
		query = query.Distinct()
	}

	return query, nil
}

func processOrderBys(query r.Term, ar *ArRethinkDb) r.Term {
	if len(ar.Query().OrderBys) > 0 {
		orderBys := []interface{}{}

		for _, orderBy := range ar.Query().OrderBys {
			switch orderBy.SortOrder {
			case goar.DESC: // descending
				orderBys = append(orderBys, r.Desc(orderBy.Key))
			default: // ascending
				orderBys = append(orderBys, r.Asc(orderBy.Key))
			}
		}

		query = query.OrderBy(orderBys...)
	}

	return query
}
