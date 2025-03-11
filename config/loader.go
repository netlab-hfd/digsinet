package config

// adapted from https://github.com/vsouza/go-gin-boilerplate

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init(mode string) {
	// if mode is empty, set the config name to default
	// otherwise, set the config name to the mode
	if mode == "" {
		viper.SetConfigName("config")
	} else {
		viper.SetConfigName(mode)
	}
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	viper.SetConfigType("yaml")
	viper.AddConfigPath(homedir + "/.digsinet-ng")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")
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
