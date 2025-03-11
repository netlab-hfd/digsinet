package config

import (
	"testing"

	"github.com/spf13/viper"
)

func TestConfig(t *testing.T) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("env")

	_, err := GetConfig()
	if err != nil {
		t.Fatalf("Unexpected Error while reading configuration from file: %s", err.Error())
	}
}
