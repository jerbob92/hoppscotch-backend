package config

import (
	"github.com/spf13/viper"
)

func LoadConfig() error {
	viper.SetConfigName("config")          // name of config file (without extension)
	viper.SetConfigType("yaml")            // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")               // optionally look for config in the working directory
	viper.AddConfigPath("/etc/api-config") // look for config in the api-config directory on the server

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
