package couchbase

import (
	"errors"
	"log"
	"os"
	"reflect"

	couchbase "github.com/couchbaselabs/go-couchbase"
	"github.com/joho/godotenv"
	. "github.com/obieq/goar"
)

type ArCouchbase struct {
	ActiveRecord
	ID string `json:"id,omitempty"`
	Timestamps
}

var (
	client *couchbase.Bucket
)

var connectOpts = func() map[string]string {
	opts := make(map[string]string)

	if envs, err := godotenv.Read(); err != nil {
		log.Fatal("Error loading .env file: ", err)
	} else {
		opts["uri"] = os.Getenv(envs["COUCHBASE_URI"])
		opts["port"] = os.Getenv(envs["COUCHBASE_PORT"])
		opts["pool"] = os.Getenv(envs["COUCHBASE_POOL"])
		opts["bucket_name"] = os.Getenv(envs["COUCHBASE_BUCKET_NAME"])
		opts["bucket_password"] = os.Getenv(envs["COUCHBASE_BUCKET_PASSWORD"])
	}

	return opts
}

func connect() *couchbase.Bucket {
	opts := connectOpts()

	if opts["uri"] == "" || opts["port"] == "" || opts["pool"] == "" || opts["bucket_name"] == "" {
		log.Fatalf("at least one required option is missing: ", opts)
	}

	// NOTE: format of connection endpoint for a password protected bucket
	//       bucket, err := couchbase.GetBucket("http://bucketname:bucketpass@myserver:8091/", "default", "bucket")
	bucketPassword := ""
	if opts["bucket_password"] != "" {
		bucketPassword = opts["bucket_name"] + ":" + opts["couchbase_buckect_password"] + "@"
	}
	endpoint := "http://" + bucketPassword + opts["uri"] + ":" + opts["port"] + "/"
	bucket, err := couchbase.GetBucket(endpoint, opts["pool"], opts["bucket_name"])

	if err != nil {
		log.Fatalf("Error getting bucket:  %v", err)
	}

	return bucket
}

func init() {
	client = connect()
}

func Client() *couchbase.Bucket {
	return client
}

func (ar *ArCouchbase) SetKey(key string) {
	ar.ID = key
}

func (ar *ArCouchbase) All(models interface{}, opts map[string]interface{}) (err error) {
	return errors.New("All method not supported by Couchbase.  Create a View instead.")
}

func (ar *ArCouchbase) Truncate() (numRowsDeleted int, err error) {
	// http://docs.couchbase.com/admin/admin/REST/rest-bucket-flush.html
	return -1, errors.New("Truncate method not yet implemented")
}

func (ar *ArCouchbase) Find(key interface{}) (interface{}, error) {
	self := ar.Self()
	modelVal := reflect.ValueOf(self).Elem()
	modelInterface := reflect.New(modelVal.Type()).Interface()

	err := client.Get(key.(string), modelInterface)

	// if there's no document found, then the modelInterface instance will be empty/nil for all properties
	// NOTE: given that the model is "generic" at this point, we need to use reflection in order to verify
	// TODO: is there a better, more efficient way to verify?  is there some way to leverage the couchbase sdk?
	if reflect.ValueOf(modelInterface).Elem().FieldByName("ID").String() == "" {
		modelInterface = nil
	}

	return modelInterface, err
}

func (ar *ArCouchbase) DbSave() error {
	// client.Set performs an upsert
	// return client.Set(ar.ID, 0, ar.Self())

	var err error
	var added bool = false

	if ar.UpdatedAt == nil {
		added, err = client.Add(ar.ID, 0, ar.Self())
		if err == nil && !added {
			err = errors.New("Insert Failed: key already exists")
		}
	} else {
		err = client.Set(ar.ID, 0, ar.Self())
	}

	return err
}

func (ar *ArCouchbase) DbDelete() (err error) {
	return client.Delete(ar.ID)
}

func (ar *ArCouchbase) DbSearch(models interface{}) (err error) {
	return errors.New("Search method not supported by Couchbase.  Create a View instead.")
}

//func processPlucks(query r.Term, ar *ArRethinkDb) r.Term {
//if plucks := ar.Query().Plucks; plucks != nil {
//query = query.Pluck(plucks...)
//}

//return query
//}
