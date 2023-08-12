package config

type Kafka struct {
	Brokers                string `mapstructure:"brokers" json:"brokers"`
	MinBytes               int    `mapstructure:"min-bytes" json:"min-bytes"`
	MaxBytes               int    `mapstructure:"max-bytes" json:"max-bytes"`
	AllowAutoTopicCreation bool   `mapstructure:"allow-auto-topic-creation" json:"allow-auto-topic-creation"`
}
