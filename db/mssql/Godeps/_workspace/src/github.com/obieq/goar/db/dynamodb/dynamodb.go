package dynamodb

import (
	"errors"
	"log"

	aws "github.com/AdRoll/goamz/aws"
	dynamo "github.com/AdRoll/goamz/dynamodb"
	goar "github.com/obieq/goar"
)

const DB_PRIMARY_KEY_NAME string = "id"
const MODEL_PRIMARY_KEY_NAME string = "ID"

type ArDynamodb struct {
	goar.ActiveRecord
	ID string `json:"id,omitempty"`
	goar.Timestamps
}

// interface assertions
// https://splice.com/blog/golang-verify-type-implements-interface-compile-time/
var _ goar.Persister = (*ArDynamodb)(nil)

var (
	clients = map[string]*dynamo.Server{}
)

func connect(connName string, env string) (s *dynamo.Server) {
	c := goar.Config
	if c == nil {
		log.Panicln("goar config cannot be nil")
	}

	connKey := env + "_dynamodb_" + connName
	m, found := c.DynamoDBs[connKey]
	if !found {
		log.Panicln("dynamodb connection not found:", connKey)
	} else if m.Region == "" {
		log.Panicln("dynamodb aws region cannot be blank")
	} else if m.AccessKey == "" {
		log.Panicln("dynamodb access key cannot be blank")
	} else if m.SecretKey == "" {
		log.Panicln("dynamodb secret key cannot be blank")
	}

	var region aws.Region
	switch m.Region {
	case "useast":
		region = aws.USEast
	case "uswest":
		region = aws.USWest
	case "uswest2":
		region = aws.USWest2
	default:
		log.Panicln("invalid region:", m.Region)
	}
	auth := aws.Auth{AccessKey: m.AccessKey, SecretKey: m.SecretKey}

	return dynamo.New(auth, region)
}

func (ar *ArDynamodb) Client() *dynamo.Server {
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

func (ar *ArDynamodb) SetKey(key string) {
	ar.ID = key
}

func (ar *ArDynamodb) All(models interface{}, opts map[string]interface{}) (err error) {
	return errors.New("All method not supported by Dynamodb.  Create a View instead.")
}

func (ar *ArDynamodb) Truncate() (numRowsDeleted int, err error) {
	return -1, errors.New("Truncate method not yet implemented")
}

func (ar *ArDynamodb) Find(id interface{}, out interface{}) error {
	tbl, dynamoKey := ar.GetTableWithPrimaryKey(id)

	// NOTE: the AdRoll sdk returns an error if the key doesn't exist
	if err := tbl.GetDocument(dynamoKey, out); err != nil {
		return err
	}

	// set the ID b/c the AdRoll sdk doen't map embedded struct properties at all TODO: follow up w/ AdRoll
	out.(goar.ActiveRecordInterfacer).SetKey(id.(string))

	return nil
}

func (ar *ArDynamodb) DbSave() error {
	tbl, key := ar.GetTableWithPrimaryKey(ar.ID)
	return tbl.PutDocument(key, ar.Self())
}

// func (ar *ArDynamodb) Patch() (bool, error) {
// 	var err error
// 	var success bool = false
// 	tbl, key := ar.GetTableWithPrimaryKey(ar.ID)
//
// 	// copy self for later use if we perorm an update
// 	// reason being that when find is called, it will update
// 	// the underlying Self() instance
// 	source := reflect.ValueOf(ar.Self()).Elem().Interface()
//
// 	// query the db to determine if we're doing an insert or an update
// 	// NOTE: due to the fact this supports PATCH updates, we need to
// 	//       get the persisted instance if one exists in order to update
// 	//       the subset of fields.  If we didn't do so, then biz rule validations
// 	//       could fail b/c of the incomplete data
// 	dbInstance, err := ar.Find(ar.ID)
//
// 	if err == nil { // instance found, so update
// 		var arr []reflect.Value
//
// 		// sync db instance and self instance
// 		e := reflect.ValueOf(dbInstance).Elem()
// 		addr := e.Addr()
// 		method := addr.MethodByName("Self")
// 		destination := method.Call(arr)[0].Interface()
//
// 		// merge patch changes into existing instance
// 		// NOTE: nil/empty values don't appear to overwrite existing values
// 		err = mergo.Merge(destination, source)
//
// 		if err == nil {
// 			// set updated at timestamp
// 			updatedAt := time.Now().UTC()
// 			ar.UpdatedAt = &updatedAt // given the use of pointers, no need to use reflection
//
// 			// run validations and update
// 			success = ar.Valid()
// 			if success == true {
// 				err = tbl.PutDocument(key, destination)
// 			}
// 		}
// 	}
//
// 	return success, err
// }

func (ar *ArDynamodb) DbDelete() (err error) {
	primary := dynamo.NewStringAttribute(DB_PRIMARY_KEY_NAME, "")
	pk := dynamo.PrimaryKey{KeyAttribute: primary}
	t := dynamo.Table{Server: ar.Client(), Name: ar.ModelName(), Key: pk}

	dynamoKey := &dynamo.Key{HashKey: ar.ID}
	return t.DeleteDocument(dynamoKey)
}

func (ar *ArDynamodb) DbSearch(models interface{}) (err error) {
	return errors.New("Search method not supported by Dynamodb.  Create a View instead.")
}

func (ar *ArDynamodb) GetTableWithPrimaryKey(key interface{}) (dynamo.Table, *dynamo.Key) {
	// primary key initialization example
	//     https://github.com/AdRoll/goamz/blob/c73835dc8fc6958baf8df8656864ee4d6d04b130/dynamodb/query_builder_test.go
	//         primary := NewStringAttribute("TestHashKey", "")
	//         secondary := NewNumericAttribute("TestRangeKey", "")
	//         key := PrimaryKey{primary, secondary}
	primary := dynamo.NewStringAttribute(DB_PRIMARY_KEY_NAME, "")
	pk := dynamo.PrimaryKey{KeyAttribute: primary}
	t := dynamo.Table{Server: ar.Client(), Name: ar.ModelName(), Key: pk}
	dynamoKey := &dynamo.Key{HashKey: key.(string)}

	return t, dynamoKey
}
