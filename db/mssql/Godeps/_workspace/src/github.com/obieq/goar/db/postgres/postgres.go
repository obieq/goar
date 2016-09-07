package postgres

import (
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/jinzhu/gorm"
	. "github.com/obieq/goar"
)

type ArPostgres struct {
	ActiveRecord
	ID int `gorm:"primary_key" json:"id,omitempty"`
	//ID string `sql:"type:varchar(36)" gorm:"primary_key" json:"id,omitempty"`
	Timestamps
}

// interface assertions
// https://splice.com/blog/golang-verify-type-implements-interface-compile-time/
var _ Persister = (*ArPostgres)(nil)
var _ RDBMSer = (*ArPostgres)(nil)

var (
	clients = map[string]gorm.DB{}
)

func connect(connName string, env string) (client gorm.DB) {
	c := Config
	if c == nil {
		log.Panic("goar config cannot be nil")
	}

	connKey := env + "_postgresql_" + connName
	m, found := c.PostgresqlDBs[connKey]
	if !found {
		log.Panic("postgresql connection not found:", connKey)
	}

	connString := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", m.Server, m.Port, m.DBName, m.Username, m.Password)

	if m.Debug {
		log.Println("connString:", connString)
	}

	db, err := gorm.Open("postgres", connString)
	if err != nil {
		log.Panic("Open postgresql database failed:", err)
	}

	// set connection properties
	db.DB().SetMaxIdleConns(m.MaxIdleConnections)
	db.DB().SetMaxOpenConns(m.MaxOpenConnections)

	// set log mode
	db.LogMode(m.Debug)

	// test the connection
	if err = db.DB().Ping(); err != nil {
		log.Panic("postgresql ping failed:", err.Error())
	}

	//return connection
	return db
}

func (ar *ArPostgres) Client() gorm.DB {
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

func (ar *ArPostgres) SetKey(key string) {
	// TODO: set guid key here once that's implemented
	//ar.ID = key
}

func (ar *ArPostgres) All(models interface{}, opts map[string]interface{}) (err error) {
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

func (ar *ArPostgres) Truncate() (numRowsDeleted int, err error) {
	client := ar.Client()
	tblName := client.NewScope(ar.Self()).TableName()
	return -1, client.Exec(`TRUNCATE "` + tblName + `";`).Error
}

func (ar *ArPostgres) Find(id interface{}, out interface{}) error {
	//result, err := client.Get(ar.ModelName(), id.(string))

	//if result != nil {
	//err = result.Value(&out)
	//} else {
	//err = errors.New("record not found")
	//}

	client := ar.Client()
	return client.First(out, id).Error
	//return nil
}

func (ar *ArPostgres) DbSave() error {
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

func (ar *ArPostgres) DbDelete() (err error) {
	//return client.Purge(ar.ModelName(), ar.ID)
	return nil
}

func (ar *ArPostgres) DbSearch(models interface{}) (err error) {
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

func (ar *ArPostgres) SpExecResultSet(spName string, params map[string]interface{}, models interface{}) (err error) {
	return errors.New("postgres.SpExecResultSet not implemented")
}

func buildSpParams(params map[string]interface{}) string {
	log.Fatalf("postgres.buildSpParams not implemented")

	return ""
}

func mapResults(results interface{}, models interface{}) (err error) {
	return nil
}

func processWhereConditions(ar *ArPostgres) (query string, err error) {
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

func processSorts(ar *ArPostgres) (sort string) {
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
