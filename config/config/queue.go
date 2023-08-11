package config

type Queue struct {
	Redis  Redis
	Kafka  Kafka
	Prefix string
}
