package queue

import (
	"gin-example/global"
	"gin-example/queue/core"
	"time"
)

type TestJob struct{}

var TestJobApp = &TestJob{}

func init() {
	core.QueueApp.Bind(TestJobApp)
}

// 队列名
func (*TestJob) GetName() string {
	return global.TestJob
}

// 连接驱动类型
func (*TestJob) GetConnType() core.ConnType {
	return core.Redis
}

// 处理消费
func (t *TestJob) Handel(message string) error {
	time.Sleep(time.Millisecond * 50)
	// p := queue.TestJobPayload{}
	/*err := sonic.UnmarshalString(message, &p)
	if err != nil {
		return fmt.Errorf("unmarshal error: %s", err)
	}*/
	// utils.DelCache(p.Key)
	return nil
}

// 是否开启消费监听
func (t *TestJob) Enable() bool {
	return true
}

func (t *TestJob) GetConsumerNumber() int {
	return 100
}

func (t *TestJob) GetRetryCount() int {
	return 3
}

func (t *TestJob) GetRateLimit() int {
	return 200
}
