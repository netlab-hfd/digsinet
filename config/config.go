package config

import "fmt"

type GnmiConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Publish  bool   `mapstructure:"publish"`
}

type KafkaBroker struct {
	Hostname string `mapstructure:"hostname"`
	Port     int    `mapstructure:"port"`
}

func (k *KafkaBroker) ConnectionString() string {
	return fmt.Sprintf("%s:%d", k.Hostname, k.Port)
}

type KafkaConfig struct {
	Brokers []KafkaBroker `mapstructure:"brokers"`
}

type RestConfig struct {
	Address    string `mapstructure:"address"`
	AuthKey    string `mapstructure:"auth.key"`
	AuthSecret string `mapstructure:"auth.secret"`
}

type Configuration struct {
	Gnmi  GnmiConfig  `mapstructure:"gnmi"`
	Kafka KafkaConfig `mapstructure:"kafka"`
	Http  RestConfig  `mapstructure:"http"`
}
