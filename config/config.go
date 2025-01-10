package config

// adapted from https://github.com/vsouza/go-gin-boilerplate

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var config *viper.Viper

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init(env string) {
	var err error
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName("default")
	config.AddConfigPath("config/")
	err = config.ReadInConfig()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to read config")
	}

	envConfig := viper.New()
	envConfig.SetConfigType("yaml")
	envConfig.AddConfigPath("config/")
	envConfig.SetConfigName(env)
	err = envConfig.ReadInConfig()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to parse env configuration file")
	}

	err = config.MergeConfigMap(envConfig.AllSettings())
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to merge config")
		return
	}
}

func GetConfig() *viper.Viper {
	return config
}
