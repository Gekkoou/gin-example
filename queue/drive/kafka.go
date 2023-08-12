package drive

import (
	"context"
	"errors"
	"fmt"
	"gin-example/config/config"
	"github.com/bytedance/sonic"
	"github.com/segmentio/kafka-go"
	"net"
	"strconv"
	"strings"
)

type _kafka struct {
	Drive
	Produce *kafka.Writer
	Comsume *kafka.Reader
	Cfg     config.Kafka

	LastMessage kafka.Message
}

func NewKafka(topic string, cfg config.Kafka, prefix, failureSuffix string) (Interface, error) {
	k := &_kafka{
		Cfg: cfg,
	}
	k.Name, k.FailureName = getQueueNames(topic, prefix, failureSuffix)
	// kafka 不能使用 冒号作为 topic , 修改为 -
	k.Name = k.changeName(k.Name)
	k.FailureName = k.changeName(k.FailureName)
	k.Produce = k.newWriter()
	k.Comsume = k.getReader()
	err := k.checkTopic()
	return k, err
}

func (*_kafka) changeName(s string) string {
	return strings.Replace(s, ":", "-", len(s))
}

func (k *_kafka) checkTopic() (err error) {
	conn, err := kafka.Dial("tcp", k.Cfg.Brokers)
	if err != nil {
		return
	}
	defer conn.Close()
	partitions, err := conn.ReadPartitions()
	if err != nil {
		return
	}
	m := map[string]struct{}{}
	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}
	nameMap := map[string]struct{}{
		k.Name:        {},
		k.FailureName: {},
	}
	for name, _ := range nameMap {
		if _, ok := m[name]; ok {
			delete(nameMap, name)
		}
	}
	if len(nameMap) > 0 {
		if k.Cfg.AllowAutoTopicCreation {
			// 创建topic
			err = k.createTopic(conn, nameMap)
		} else {
			nameString, _ := sonic.MarshalString(nameMap)
			err = errors.New(fmt.Sprintf("kafka 主题未创建且自动创建设置为false:%s", nameString))
		}
	}
	return
}
func (k *_kafka) createTopic(conn *kafka.Conn, m map[string]struct{}) (err error) {
	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	var topicConfigs []kafka.TopicConfig
	for name, _ := range m {
		tc := kafka.TopicConfig{
			Topic:             name,
			NumPartitions:     k.Cfg.NumPartitions,     // 分区数量
			ReplicationFactor: k.Cfg.ReplicationFactor, // 副本数量，不能超过主机的数量，否则创建会失败
		}
		topicConfigs = append(topicConfigs, tc)
	}
	err = controllerConn.CreateTopics(topicConfigs...)
	return
}

// 生产者
func (k *_kafka) newWriter() *kafka.Writer {
	brokers := strings.Split(k.Cfg.Brokers, ",")
	return &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: k.Cfg.AllowAutoTopicCreation, // 是否允许不存在的 topic 自动创建, 需要注意，如果是自动创建的话，创建完之后集群会进行选举，不要马上发消息
	}
}

// 消费者
func (k *_kafka) getReader() *kafka.Reader {
	brokers := strings.Split(k.Cfg.Brokers, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  k.Name,
		Topic:    k.Name,
		MinBytes: k.Cfg.MinBytes,
		MaxBytes: k.Cfg.MaxBytes,
	})
}

func (k *_kafka) Push(ctx context.Context, message string) (err error) {
	err = k.Produce.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(""),
			Topic: k.Name,
			Value: []byte(message),
		},
	)
	return
}

func (k *_kafka) PushFailure(ctx context.Context, message string) (err error) {
	err = k.Produce.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(""),
			Topic: k.FailureName,
			Value: []byte(message),
		},
	)
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

func (k *_kafka) GetMessage(ctx context.Context) (string, error) {
	message, err := k.Comsume.FetchMessage(ctx)
	if err != nil {
		return "", err
	}
	k.LastMessage = message
	return string(message.Value), nil
}

func (k *_kafka) CommitMessage(ctx context.Context) (err error) {
	err = k.Comsume.CommitMessages(ctx, k.LastMessage)
	return
}
