package queue

import (
	"fmt"
	"gin-example/queue/core"
	"gin-example/utils"
	"github.com/bytedance/sonic"
)

type DelCacheTask struct{}

var DelCacheTaskApp = &DelCacheTask{}

func init() {
	core.QueueApp.Bind(DelCacheTaskApp)
}

// 队列名
func (j *DelCacheTask) GetName() string {
	return "del-cache"
}

// 连接驱动类型
func (j *DelCacheTask) GetConnType() core.ConnType {
	return core.Kafka
}

// 处理消费
func (t *DelCacheTask) Handel(message string) error {
	p := DelCachePayload{}
	err := sonic.UnmarshalString(message, &p)
	if err != nil {
		return fmt.Errorf("unmarshal error: %s", err)
	}
	utils.DelCache(p.Key)
	return nil
}

// 是否开启消费监听
func (t *DelCacheTask) Enable() bool {
	return true
}

func (t *DelCacheTask) GetConsumerNumber() int {
	return 3
}

type DelCachePayload struct {
	Key string
}
