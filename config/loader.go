package config

// adapted from https://github.com/vsouza/go-gin-boilerplate

import (
	"github.com/spf13/viper"
)

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init(env string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("env")
}

func GetConfig() (*Configuration, error) {
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var cfg Configuration
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
