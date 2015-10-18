package gas

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
}

func (c *Config) Load(appName string, configFilename string, recurse bool) error {
	viper.Reset()
	viper.SetConfigName(strings.Split(configFilename, ".")[0])

	if recurse {
		loadConfigFileRecursively(configFilename)
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Panicln("Fatal error reading "+appName+" config file:", err)
	}
	return err
}

func loadConfigFileRecursively(configFilename string) error {
	wd, err := os.Getwd()
	if err != nil {
		log.Panicln("An error occurred in gas while trying to get the calling app's working directory:", err)
	}

	arr := strings.Split(wd+"/", "/")
	i := len(arr) - 1
	for i > 1 {
		fp := strings.Join(arr[0:i], "/") + "/"
		log.Println("attempting to load config file from:", fp)
		// see if config file exists for the current path
		if _, err := os.Stat(fp + configFilename); err == nil {
			log.Println("found "+configFilename+" in:", fp)
			viper.AddConfigPath(fp)
			break
		}
		i--
	}
	return err
}

func GetString(path string) string {
	v := viper.GetString(path)
	if strings.HasPrefix(v, "ENV[") {
		return getEnv(v[4:len(v)-1], path)
	}
	return v
}

func GetInt(path string) int {
	v := viper.GetString(path)
	if strings.HasPrefix(v, "ENV[") {
		i, err := strconv.Atoi(getEnv(v[4:len(v)-1], path))
		if err != nil {
			log.Panicln("cannot cast env value to int for:", path)
		}
		return i
	}
	return viper.GetInt(path)
}

func GetBool(path string) bool {
	v := viper.GetString(path)
	if strings.HasPrefix(v, "ENV[") {
		b, err := strconv.ParseBool(getEnv(v[4:len(v)-1], path))
		if err != nil {
			log.Panicln("cannot cast env value to int for:", path)
		}
		return b
	}
	return viper.GetBool(path)
}

func getEnv(name string, path string) string {
	v := os.Getenv(name)
	if v == "" {
		log.Panicln("env value cannot be blank:", path)
	}
	return v
}
