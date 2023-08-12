package config

type Queue struct {
	Redis         Redis
	Kafka         Kafka
	Prefix        string `mapstructure:"prefix" json:"prefix"`
	FailureSuffix string `mapstructure:"failure-suffix" json:"failure-suffix"`
}
