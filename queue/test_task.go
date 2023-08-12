package queue

import (
	"fmt"
	"gin-example/queue/core"
)

type TestTask struct{}

var TestTaskApp = &TestTask{}

func init() {
	core.QueueApp.Bind(TestTaskApp)
}

// 队列名
func (j *TestTask) GetName() string {
	return "test-job"
}

// 连接驱动类型
func (j *TestTask) GetConnType() core.ConnType {
	return core.Kafka
}

// 处理消费
func (t *TestTask) Handel(message string) error {
	fmt.Println("test-job的消费：", message)
	return nil
}

// 是否开启消费监听
func (t *TestTask) Enable() bool {
	return true
}

func (t *TestTask) GetConsumerNumber() int {
	return 3
}

func (t *TestTask) GetRetryCount() int {
	return 3
}
