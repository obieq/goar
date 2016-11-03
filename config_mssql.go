package goar

import (
	"github.com/obieq/gas"
	"github.com/spf13/viper"
)

// MSSQLConfig => contains mssql db connection info
type MSSQLConfig struct {
	ConnectionName     string
	Server             string
	FailoverPartner    string
	FailoverPort       int
	Port               int
	DBName             string
	Username           string
	Password           string
	MaxIdleConnections int
	MaxOpenConnections int
	Debug              bool
}

func (c *config) loadMSSQL() {
	for connKey := range viper.GetStringMap(c.Environment + "." + MSSQL) {
		connName := c.Environment + "_" + MSSQL + "_" + connKey
		path := c.Environment + "." + MSSQL + "." + connKey + "."
		mssql := &MSSQLConfig{}
		mssql.ConnectionName = connName

		mssql.Server = gas.GetString(path + "server")
		mssql.Port = gas.GetInt(path + "port")

		mssql.FailoverPartner = gas.GetString(path + "failoverpartner")
		mssql.FailoverPort = gas.GetInt(path + "failoverport")

		mssql.DBName = gas.GetString(path + "dbname")
		mssql.Username = gas.GetString(path + "username")
		mssql.Password = gas.GetString(path + "password")
		mssql.MaxIdleConnections = gas.GetInt(path + "maxidleconnections")
		mssql.MaxOpenConnections = gas.GetInt(path + "maxopenconnections")
		mssql.Debug = gas.GetBool(path + "debug")

		c.MSSQLDBs[connName] = mssql
	}
}
