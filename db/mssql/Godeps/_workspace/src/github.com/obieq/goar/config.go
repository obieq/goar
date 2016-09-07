package goar

import (
	"log"

	"github.com/obieq/gas"
)

const MSSQL = "mssql"
const RETHINKDB = "rethinkdb"
const POSTGRESQL = "postgresql"
const DYNAMODB = "dynamodb"
const ORCHESTRATE = "orchestrate"
const COUCHBASE = "couchbase"

var Config *config

func init() {
	log.Println("calling goar.config.init()")
	Config = newConfig()
}

// Config => contains arrays of connections for a given db provider
type config struct {
	gas.Config
	Environment    string
	MSSQLDBs       map[string]*MSSQLConfig
	RethinkDBs     map[string]*RethinkDBConfig
	PostgresqlDBs  map[string]*PostgresqlDBConfig
	DynamoDBs      map[string]*DynamoDBConfig
	OrchestrateDBs map[string]*OrchestrateConfig
	CouchbaseDBs   map[string]*CouchbaseConfig
}

func newConfig() *config {
	c := &config{}
	c.MSSQLDBs = map[string]*MSSQLConfig{}
	c.RethinkDBs = map[string]*RethinkDBConfig{}
	c.PostgresqlDBs = map[string]*PostgresqlDBConfig{}
	c.DynamoDBs = map[string]*DynamoDBConfig{}
	c.OrchestrateDBs = map[string]*OrchestrateConfig{}
	c.CouchbaseDBs = map[string]*CouchbaseConfig{}

	err := c.Load("goar", "config.json", true)

	// [99858326] could create a race condition, so don't call AutomaticEnv()
	// viper.AutomaticEnv() // loads env

	// parse Environment
	c.Environment = gas.GetString("environment")
	if c.Environment == "" {
		log.Panicln("goar config environment cannot be blank:", err)
	}

	// parse mssql connections
	c.loadMSSQL()

	// parse rethinkdb connections
	c.loadRethinkDB()

	// parse postgresql connections
	c.loadPostgresql()

	// parse dynamodb connections
	c.loadDynamoDB()

	// parse orchestrate connections
	c.loadOrchestrate()

	// parse couchbase connections
	c.loadCouchbase()

	return c
}
