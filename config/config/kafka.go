package config

type Kafka struct {
	Brokers                string
	MinBytes               int
	MaxBytes               int
	AllowAutoTopicCreation bool
}
