package goar

import (
	"github.com/obieq/gas"
	"github.com/spf13/viper"
)

// CouchbaseConfig => contains couchbase db connection info
type CouchbaseConfig struct {
	ConnectionName string
	ClusterAddress string
	BucketName     string
	BucketPassword string
}

func (c *config) loadCouchbase() {
	for connKey := range viper.GetStringMap(c.Environment + "." + COUCHBASE) {
		connName := c.Environment + "_" + COUCHBASE + "_" + connKey
		path := c.Environment + "." + COUCHBASE + "." + connKey + "."
		couchbase := &CouchbaseConfig{}
		couchbase.ConnectionName = connName
		couchbase.ClusterAddress = gas.GetString(path + "cluster_address")
		couchbase.BucketName = gas.GetString(path + "bucket_name")
		couchbase.BucketPassword = gas.GetString(path + "bucket_password")

		c.CouchbaseDBs[connName] = couchbase
	}
}
