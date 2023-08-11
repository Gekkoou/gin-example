package drive

import (
	"context"
	"gin-example/config/config"
	"github.com/segmentio/kafka-go"
	"log"
	"strings"
)

type _kafka struct {
	Produce     *kafka.Writer
	Comsume     *kafka.Reader
	Conf        config.Kafka
	Name        string
	Prefix      string
	LastMessage kafka.Message
}

// 生产者
func (k *_kafka) newKafkaWriter() *kafka.Writer {
	brokers := strings.Split(k.Conf.Brokers, ",")
	return &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  k.Name,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: k.Conf.AllowAutoTopicCreation, // 是否允许不存在的topic 自动创建 , 需要注意，如果是自动创建的话，创建完之后集群会进行选举，不要马上发消息
	}
}

func NewKafka(topic string, conf config.Kafka, prefix string) Interface {
	k := &_kafka{
		Conf: conf,
	}
	k.Name = topic
	k.Prefix = prefix
	k.Produce = k.newKafkaWriter()
	k.Comsume = k.getKafkaReader()
	return k
}

// 消费者
func (k *_kafka) getKafkaReader() *kafka.Reader {
	brokers := strings.Split(k.Conf.Brokers, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  k.Name,
		Topic:    k.Name,
		MinBytes: k.Conf.MinBytes,
		MaxBytes: k.Conf.MaxBytes,
	})
}

func (k *_kafka) Push(ctx context.Context, message string) (err error) {
	err = k.Produce.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(""),
			Value: []byte(message),
		},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
	return
}

/*// 隐式提交 会有丢数据风险 , 不使用这个
func (k *_kafka) GetMessage_() (string, error) {
	message, err := k.Comsume.ReadMessage(context.Background())
	if err != nil {
		return "", err
	}
	return string(message.Value), nil
}*/

func (k *_kafka) CommitMessage(ctx context.Context) (err error) {
	err = k.Comsume.CommitMessages(ctx, k.LastMessage)
	return
}

func (k *_kafka) GetMessage(ctx context.Context) (string, error) {
	message, err := k.Comsume.FetchMessage(ctx)
	if err != nil {
		return "", err
	}
	k.LastMessage = message
	return string(message.Value), nil
}
