package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func LoadConfig() error {
	viper.SetConfigName("config")          // name of config file (without extension)
	viper.SetConfigType("yaml")            // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")               // optionally look for config in the working directory
	viper.AddConfigPath("/etc/api-config") // look for config in the api-config directory on the server
	configPath := os.Getenv("CONFIG_PATH_DIR")
	if configPath != "" {
		viper.AddConfigPath(configPath)
	}
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// Check for deprecated fields
	if viper.GetString("database.address") != "" {
		log.Warningln("database.address is deprecated and will be removed in future releases, please use database.host, database.port and database.driver fields for database connectivity")
	}
	return nil
}
