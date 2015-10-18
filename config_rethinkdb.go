package goar

import (
	"strings"

	"github.com/obieq/gas"
	"github.com/spf13/viper"
)

// RethinkDBConfig => contains rethinkdb db connection info
type RethinkDBConfig struct {
	ConnectionName     string
	Addresses          []string
	DBName             string
	AuthKey            string
	DiscoverHosts      bool
	MaxIdleConnections int
	MaxOpenConnections int
	Debug              bool
}

// "addresses": "ec2-52-7-204-235.compute-1.amazonaws.com:28015",
// "dbname": "VINISO",
// "authkey": "",
// "maxidleconnections": 10,
// "maxopenconnections": 100,
// "debug": true

func (c *config) loadRethinkDB() {
	for connKey := range viper.GetStringMap(c.Environment + "." + RETHINKDB) {
		connName := c.Environment + "_" + RETHINKDB + "_" + connKey
		path := c.Environment + "." + RETHINKDB + "." + connKey + "."
		rethinkdb := &RethinkDBConfig{}
		rethinkdb.ConnectionName = connName
		rethinkdb.Addresses = parseRethinkDbAddresses(gas.GetString(path + "addresses"))
		rethinkdb.DBName = gas.GetString(path + "dbname")
		rethinkdb.AuthKey = gas.GetString(path + "authkey")
		rethinkdb.DiscoverHosts = gas.GetBool(path + "discoverhosts")
		rethinkdb.MaxIdleConnections = gas.GetInt(path + "maxidleconnections")
		rethinkdb.MaxOpenConnections = gas.GetInt(path + "maxopenconnections")
		rethinkdb.Debug = gas.GetBool(path + "debug")

		c.RethinkDBs[connName] = rethinkdb
	}
}

func parseRethinkDbAddresses(configValue string) []string {
	addresses := []string{}
	for _, address := range strings.Split(configValue, ",") {
		addresses = append(addresses, address)
	}

	return addresses
}
