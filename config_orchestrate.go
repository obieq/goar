package goar

import (
	"github.com/obieq/gas"
	"github.com/spf13/viper"
)

// OrchestrateConfig => contains orchestrate db connection info
type OrchestrateConfig struct {
	ConnectionName string
	APIKey         string
}

func (c *config) loadOrchestrate() {
	for connKey := range viper.GetStringMap(c.Environment + "." + ORCHESTRATE) {
		connName := c.Environment + "_" + ORCHESTRATE + "_" + connKey
		path := c.Environment + "." + ORCHESTRATE + "." + connKey + "."
		orchestrate := &OrchestrateConfig{}
		orchestrate.ConnectionName = connName
		orchestrate.APIKey = gas.GetString(path + "apikey")

		c.OrchestrateDBs[connName] = orchestrate
	}
}
