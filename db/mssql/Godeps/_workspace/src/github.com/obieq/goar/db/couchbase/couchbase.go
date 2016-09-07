package couchbase

import (
	"errors"
	"log"

	gocb "github.com/couchbase/gocb"
	goar "github.com/obieq/goar"
	"github.com/satori/go.uuid"
)

type ArCouchbase struct {
	goar.ActiveRecord
	ID string `json:"id,omitempty"`
	goar.Timestamps
}

// interface assertions
// https://splice.com/blog/golang-verify-type-implements-interface-compile-time/
var _ goar.Persister = (*ArCouchbase)(nil)

var (
	clients = map[string]*gocb.Bucket{}
)

func connect(connName string, env string) *gocb.Bucket {
	cfg := goar.Config
	if cfg == nil {
		log.Panicln("goar couchbase config cannot be nil")
	}

	connKey := env + "_couchbase_" + connName
	m, found := cfg.CouchbaseDBs[connKey]
	if !found {
		log.Panicln("couchbase db connection not found:", connKey)
	} else if m.ClusterAddress == "" {
		log.Panicln("couchbase cluster address cannot be blank")
	} else if m.BucketName == "" {
		log.Panicln("couchbase bucket name cannot be blank")
	} else if m.BucketPassword == "" {
		log.Println("---- WARNING --- bucket password is blank")
	}

	cluster, _ := gocb.Connect("couchbase://" + m.ClusterAddress + "/")
	bucket, err := cluster.OpenBucket(m.BucketName, m.BucketPassword)

	if err != nil {
		log.Fatalf("Error getting bucket:  %v", err)
	}

	// bucket.SetTranscoder(TestTranscoder{})

	return bucket
}

func (ar *ArCouchbase) Client() *gocb.Bucket {
	self := ar.Self()
	connectionKey := self.DBConnectionName() + "_" + self.DBConnectionEnvironment()
	if self == nil {
		log.Panic("couchbase ar.Self() cannot be blank!")
	}

	conn, found := clients[connectionKey]
	if !found {
		conn = connect(self.DBConnectionName(), self.DBConnectionEnvironment())
		clients[connectionKey] = conn
	}

	return conn
}

func (ar *ArCouchbase) SetKey(key string) {
	ar.ID = key
}

func (ar *ArCouchbase) All(models interface{}, opts map[string]interface{}) (err error) {
	return errors.New("All method not supported by Couchbase.  Create a View instead.")
}

func (ar *ArCouchbase) Truncate() (numRowsDeleted int, err error) {
	// http://docs.couchbase.com/admin/admin/REST/rest-bucket-flush.html
	err = ar.Client().Manager("user-name", "password").Flush()
	log.Fatal("couchbase.Truncate() failed with error: ", err)
	return -1, err
}

func (ar *ArCouchbase) Find(id interface{}, out interface{}) error {
	_, err := ar.Client().Get(id.(string), &out)
	return err
}

func (ar *ArCouchbase) DbSave() error {
	var err error
	var cas gocb.Cas

	if ar.UpdatedAt == nil {
		if ar.ID == "" { // auto-generate the document id if the client didn't provide one
			ar.ID = uuid.NewV4().String()
		}
		cas, err = ar.Client().Insert(ar.ID, ar.Self(), 0)
		if err == nil && cas == 0 {
			err = errors.New("Insert Failed: key already exists")
		}
	} else {
		// err = client.Set(ar.ID, 0, ar.Self())
		cas, err = ar.Client().Replace(ar.ID, ar.Self(), 0, 0)
	}

	return err
}

func (ar *ArCouchbase) DbDelete() (err error) {
	_, err = ar.Client().Remove(ar.ID, 0)
	return err
}

func (ar *ArCouchbase) DbSearch(models interface{}) (err error) {
	return errors.New("Search method not supported by Couchbase.  Create a View instead.")
}

func (ar *ArCouchbase) N1qlQuery(query string, models *[]interface{}) (err error) {
	var rows gocb.ViewResults
	n1qlQuery := gocb.NewN1qlQuery(query)
	if rows, err = ar.Client().ExecuteN1qlQuery(n1qlQuery, nil); err != nil {
		return err
	}

	var row interface{}
	for rows.Next(&row) {
		*models = append(*models, row)
	}
	if err := rows.Close(); err != nil {
		log.Println("N1QL query error: %s\n", err)
		return err
	}

	return err
}

//func processPlucks(query r.Term, ar *ArRethinkDb) r.Term {
//if plucks := ar.Query().Plucks; plucks != nil {
//query = query.Pluck(plucks...)
//}

//return query
//}
