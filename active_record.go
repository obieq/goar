package goar

import (
	"log"
	"reflect"
	"strings"
	"time"

	gas "github.com/obieq/gas"
	validations "github.com/obieq/goar-validations"
)

var modelNames = map[string]string{}

type Validater interface {
	Valid() bool
	Validate()
	Errors() map[string]*validations.ValidationError
}

type Persister interface {
	DbSave() (err error)
	DbDelete() (err error)
	DbSearch(results interface{}) error
}

type RDBMSer interface {
	SpExecResultSet(spName string, params map[string]interface{}, results interface{}) error
}

type ActiveRecordInterfacer interface {
	Validater
	Querier
	ModelName() string
	DBConnectionName() string        // EX: aws1, aws2, azure1, azure2, default
	DBConnectionEnvironment() string // EX: dev, qa, ci, prod
	SetKey(string)
	//PrimaryKey() string
	Self() ActiveRecordInterfacer
	SetSelf(ActiveRecordInterfacer)
	//Query() *Query
	//SetQuery(*Query)
	Truncate() (numRowsDeleted int, err error)
	All(results interface{}, opts map[string]interface{}) error
	Find(id interface{}, out interface{}) error
	Save() (success bool, err error)
	Delete() error
}

type CustomModelNamer interface {
	CustomModelName() string
}

type Timestamps struct {
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `xorm:"updated" json:"updated_at,omitempty"`
}

type ActiveRecord struct {
	validations.Validation
	self  ActiveRecordInterfacer
	query *Query
}

func (ar *ActiveRecord) ModelName() string {
	var name, key string

	if v, ok := ar.self.(CustomModelNamer); ok {
		name = v.CustomModelName()
	} else {
		key = reflect.TypeOf(ar.self).String()
		name = modelNames[key]
	}

	if name == "" {
		arr := strings.Split(key, ".")
		structName := arr[len(arr)-1]
		name = gas.String(gas.String(structName).Pluralize()).Underscore()
		modelNames[key] = name
	}

	return name
}

//func (ar *ActiveRecord) PrimaryKey() string {
//return "Id"
//}

func ToAR(ari ActiveRecordInterfacer) ActiveRecordInterfacer {
	ari.SetSelf(ari)
	ari.SetQuery(NewQuery())
	return ari
}

func (ar *ActiveRecord) Self() ActiveRecordInterfacer {
	return ar.self
}

func (ar *ActiveRecord) SetSelf(ari ActiveRecordInterfacer) {
	ar.self = ari
}

func (ar *ActiveRecord) Query() *Query {
	return ar.query
}

func (ar *ActiveRecord) SetQuery(query *Query) {
	ar.query = query
}

func (ar *ActiveRecord) Pluck(keys ...interface{}) *ActiveRecord {
	ar.Query().Plucks = keys
	return ar
}

func (ar *ActiveRecord) Where(where QueryCondition) *ActiveRecord {
	ar.Query().WhereConditions = append(ar.Query().WhereConditions, where)
	return ar
}

//func (ar *ActiveRecord) Or(or QueryCondition) *ActiveRecord {
//ar.Query().OrConditions = append(ar.Query().WhereConditions, or)
//return ar
//}

func (ar *ActiveRecord) Sum(fields ...interface{}) *ActiveRecord {
	ar.Query().Aggregations[SUM] = fields

	return ar
}

func (ar *ActiveRecord) Distinct() *ActiveRecord {
	ar.Query().Distinct = true

	return ar
}

func (ar *ActiveRecord) Order(orderBy OrderBy) *ActiveRecord {
	ar.Query().OrderBys = append(ar.Query().OrderBys, orderBy)
	return ar
}

func (ar *ActiveRecord) Run(results interface{}) error {
	err := ar.Self().(Persister).DbSearch(results)
	if err == nil {
		// reset the query struct for future queries
		ar.SetQuery(NewQuery())
	}

	return err
}

func (ar *ActiveRecord) Valid() bool {
	ar.self.Validate()
	return !ar.Validation.HasErrors()
}

func (ar *ActiveRecord) Errors() map[string]*validations.ValidationError {
	//ar.self.Validate() // TODO: is this call necessary???
	return ar.Validation.ErrorMap()
}

func (ar *ActiveRecord) Save() (success bool, err error) {
	e := reflect.ValueOf(ar.Self()).Elem()
	if err = Callback("BeforeSave", e.Addr(), nil); err != nil {
		return false, err
	}

	if ar.Valid() {
		// set timestamps
		//  1) CreatedAt is set upon create
		//     NOTE: UpdatedAt is nil
		//  2) UpdatedAt is set upon subsequent updates
		t := time.Now().UTC()
		f := e.FieldByName("CreatedAt")
		// log.Println("**********************************", reflect.TypeOf(f.Interface()).String())
		if f.IsValid() {
			if reflect.TypeOf(f.Interface()).String() == "time.Time" {
				if f.Interface().(time.Time).IsZero() {
					f.Set(reflect.ValueOf(t))
				} else {
					f = e.FieldByName("UpdatedAt")
					f.Set(reflect.ValueOf(t))
				}
			} else {
				if f.IsNil() {
					f.Set(reflect.ValueOf(&t))
				} else {
					f = e.FieldByName("UpdatedAt")
					f.Set(reflect.ValueOf(&t))
				}
			}
		}

		// if f.IsValid() {
		// 	log.Println("CreatedAt value:", f.Interface())
		// 	log.Println("CreatedAt type:", reflect.TypeOf(f.Interface()))
		// 	log.Println(f.Interface())
		// 	log.Println(reflect.Zero(reflect.TypeOf(f.Interface())))
		// 	log.Println(reflect.Zero(reflect.TypeOf(f.Interface())).Interface())
		// 	// if f.Interface() == reflect.Zero(reflect.TypeOf(f.Interface())) {
		// 	// if f.Interface() == reflect.Zero(reflect.TypeOf(f)) {
		// 	if f.Interface() == reflect.Zero(reflect.TypeOf(f.Interface())).Interface() {
		// 		log.Println("Obie")
		// 		f.Set(reflect.ValueOf(t))
		// 	} else {
		// 		log.Println("Gigi")
		// 		f = e.FieldByName("UpdatedAt")
		// 		f.Set(reflect.ValueOf(t))
		// 	}
		// }

		// save changes
		err = ar.self.(Persister).DbSave()

		// error handling
		if err == nil {
			if afterSaveErr := Callback("AfterSave", e.Addr(), nil); afterSaveErr != nil {
				log.Println(afterSaveErr) // don't return the error at this point b/c the db operation was successful
			}
		}
	}

	return !ar.Validation.HasErrors() && err == nil, err
}
func (ar *ActiveRecord) Delete() error {
	return ar.self.(Persister).DbDelete()
}

func Callback(name string, eptr reflect.Value, arg []reflect.Value) error {
	hook := eptr.MethodByName(name)
	if hook.IsValid() {
		ret := hook.Call(arg)
		if len(ret) > 0 && !ret[0].IsNil() {
			return ret[0].Interface().(error)
		}
	}

	return nil
}
