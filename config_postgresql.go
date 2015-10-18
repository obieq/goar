package goar

import (
	"github.com/obieq/gas"
	"github.com/spf13/viper"
)

// PostgresqlDBConfig => contains postgresql db connection info
type PostgresqlDBConfig struct {
	ConnectionName     string
	Server             string
	Port               int
	DBName             string
	Username           string
	Password           string
	MaxIdleConnections int
	MaxOpenConnections int
	Debug              bool
}

func (c *config) loadPostgresql() {
	for connKey := range viper.GetStringMap(c.Environment + "." + POSTGRESQL) {
		connName := c.Environment + "_" + POSTGRESQL + "_" + connKey
		path := c.Environment + "." + POSTGRESQL + "." + connKey + "."
		postgresql := &PostgresqlDBConfig{}
		postgresql.ConnectionName = connName
		postgresql.Server = gas.GetString(path + "server")
		postgresql.Port = gas.GetInt(path + "port")
		postgresql.DBName = gas.GetString(path + "dbname")
		postgresql.Username = gas.GetString(path + "username")
		postgresql.Password = gas.GetString(path + "password")
		postgresql.MaxIdleConnections = gas.GetInt(path + "maxidleconnections")
		postgresql.MaxOpenConnections = gas.GetInt(path + "maxopenconnections")
		postgresql.Debug = gas.GetBool(path + "debug")

		c.PostgresqlDBs[connName] = postgresql
	}
}
