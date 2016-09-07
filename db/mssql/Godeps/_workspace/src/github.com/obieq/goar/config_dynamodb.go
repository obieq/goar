package goar

import (
	"github.com/obieq/gas"
	"github.com/spf13/viper"
)

// DynamoDBConfig => contains dynamodb db connection info
type DynamoDBConfig struct {
	ConnectionName string
	Region         string
	AccessKey      string
	SecretKey      string
}

func (c *config) loadDynamoDB() {
	for connKey := range viper.GetStringMap(c.Environment + "." + DYNAMODB) {
		connName := c.Environment + "_" + DYNAMODB + "_" + connKey
		path := c.Environment + "." + DYNAMODB + "." + connKey + "."
		dynamodb := &DynamoDBConfig{}
		dynamodb.ConnectionName = connName
		dynamodb.Region = gas.GetString(path + "region")
		dynamodb.AccessKey = gas.GetString(path + "accesskey")
		dynamodb.SecretKey = gas.GetString(path + "secretkey")

		c.DynamoDBs[connName] = dynamodb
	}
}
