package orchestrate

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	goar "github.com/obieq/goar"
	c "github.com/orchestrate-io/gorc"
)

type ArOrchestrate struct {
	goar.ActiveRecord
	ID string `json:"id,omitempty"`
	goar.Timestamps
}

// interface assertions
// https://splice.com/blog/golang-verify-type-implements-interface-compile-time/
var _ goar.Persister = (*ArOrchestrate)(nil)

var (
	clients = map[string]*c.Client{}
)

var connectOpts = func() map[string]string {
	opts := make(map[string]string)
	opts["api_key"] = "7e839ed5-536c-405a-9d78-e18d0bbf6080"

	return opts
}

func connect(connName string, env string) (client *c.Client) {
	cfg := goar.Config
	if cfg == nil {
		log.Panicln("goar orchestrate config cannot be nil")
	}

	connKey := env + "_orchestrate_" + connName
	m, found := cfg.OrchestrateDBs[connKey]
	if !found {
		log.Panicln("orchestrate db connection not found:", connKey)
	} else if m.APIKey == "" {
		log.Panicln("orchestrate api key cannot be blank")
	}

	return c.NewClient(m.APIKey)
}

func (ar *ArOrchestrate) Client() *c.Client {
	self := ar.Self()
	connectionKey := self.DBConnectionName() + "_" + self.DBConnectionEnvironment()
	if self == nil {
		log.Panic("orchestrate ar.Self() cannot be blank!")
	}

	conn, found := clients[connectionKey]
	if !found {
		conn = connect(self.DBConnectionName(), self.DBConnectionEnvironment())
		clients[connectionKey] = conn
	}

	return conn
}

func (ar *ArOrchestrate) SetKey(key string) {
	ar.ID = key
}

func (ar *ArOrchestrate) All(models interface{}, opts map[string]interface{}) (err error) {
	var limit int = 10 // per Orchestrate's documentation: 10 default, 100 max
	var response *c.KVResults

	// set limit
	if opts["limit"] != nil {
		limit = opts["limit"].(int)
		if limit > 100 { // max limit is 100
			return errors.New("limit must be less than 100")
		}
	}

	// parse options to determine which query to use
	if opts["afterKey"] != nil {
		response, err = ar.Client().ListAfter(ar.ModelName(), opts["afterKey"].(string), limit)
	} else if opts["startKey"] != nil {
		response, err = ar.Client().ListStart(ar.ModelName(), opts["startKey"].(string), limit)
	} else {
		response, err = ar.Client().List(ar.ModelName(), limit)
	}

	if err != nil {
		return err
	}

	return mapResults(response.Results, models)
}

func (ar *ArOrchestrate) Truncate() (numRowsDeleted int, err error) {
	err = ar.Client().DeleteCollection(ar.ModelName())

	return -1, err
}

func (ar *ArOrchestrate) Find(id interface{}, out interface{}) error {
	result, err := ar.Client().Get(ar.ModelName(), id.(string))

	if result != nil {
		err = result.Value(&out)
	} else {
		err = errors.New("record not found")
	}

	return err
}

func (ar *ArOrchestrate) DbSave() error {
	var err error

	if ar.UpdatedAt != nil {
		_, err = ar.Client().Put(ar.ModelName(), ar.ID, ar.Self())
	} else {
		_, err = ar.Client().PutIfAbsent(ar.ModelName(), ar.ID, ar.Self())
	}

	return err
}

func (ar *ArOrchestrate) DbDelete() (err error) {
	return ar.Client().Purge(ar.ModelName(), ar.ID)
}

func (ar *ArOrchestrate) DbSearch(models interface{}) (err error) {
	var query, sort string
	var response *c.SearchResults
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
		if response, err = ar.Client().Search(ar.ModelName(), query, 100, 0); err != nil {
			return err
		}
	} else {
		if response, err = ar.Client().SearchSorted(ar.ModelName(), query, sort, 100, 0); err != nil {
			return err
		}
	}

	return mapResults(response.Results, models)
}

//func processPlucks(query r.Term, ar *ArRethinkDb) r.Term {
//if plucks := ar.Query().Plucks; plucks != nil {
//query = query.Pluck(plucks...)
//}

//return query
//}

func mapResults(orchestrateResults interface{}, models interface{}) (err error) {
	// now, map orchstrate's raw json to the desired active record type
	modelsv := reflect.ValueOf(models)
	if modelsv.Kind() != reflect.Ptr || modelsv.Elem().Kind() != reflect.Slice {
		panic("models argument must be a slice address")
	}
	slicev := modelsv.Elem()
	elemt := slicev.Type().Elem()

	switch t := orchestrateResults.(type) {
	case []c.KVResult:
		for _, result := range t {
			elemp := reflect.New(elemt)
			if err = result.Value(elemp.Interface()); err != nil {
				return err
			}

			slicev = reflect.Append(slicev, elemp.Elem())
		}
	case []c.SearchResult:
		for _, result := range t {
			elemp := reflect.New(elemt)
			if err = result.Value(elemp.Interface()); err != nil {
				return err
			}

			slicev = reflect.Append(slicev, elemp.Elem())
		}
	default:
		return errors.New(fmt.Sprintf("Orchestrate Response Type Not Mapped: %v", t))
	}

	// assign mapped results to the caller's supplied array
	modelsv.Elem().Set(slicev)

	return err
}

func processWhereConditions(ar *ArOrchestrate) (query string, err error) {
	var whereStmt, whereCondition string

	if len(ar.Query().WhereConditions) > 0 {
		for index, where := range ar.Query().WhereConditions {
			switch where.RelationalOperator {
			case goar.EQ: // equal
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
			case goar.GTE: // greater than or equal
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
				case goar.AND:
					whereStmt = whereStmt + " AND " + whereCondition
					//whereStmt = whereStmt.And(whereCondition)
				case goar.OR:
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

func processSorts(ar *ArOrchestrate) (sort string) {
	if len(ar.Query().OrderBys) > 0 {
		sort = ""

		for i, orderBy := range ar.Query().OrderBys {
			if i > 0 {
				sort += ","
			}

			sort += "value." + orderBy.Key + ":"

			switch orderBy.SortOrder {
			case goar.DESC: // descending
				sort += "desc"
			default: // ascending
				sort += "asc"
			}
		}
	}

	return sort
}
