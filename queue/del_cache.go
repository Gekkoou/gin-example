package queue

import (
	"fmt"
	"gin-example/global"
	"gin-example/model/queue"
	"gin-example/queue/core"
	"gin-example/utils"
	"github.com/bytedance/sonic"
	"time"
)

type DelCacheTask struct{}

var DelCacheTaskApp = &DelCacheTask{}

func init() {
	core.QueueApp.Bind(DelCacheTaskApp)
}

// 队列名
func (*DelCacheTask) GetName() string {
	return global.QueueDelCache
}

// 连接驱动类型
func (*DelCacheTask) GetConnType() core.ConnType {
	return core.Redis
}

// 处理消费
func (t *DelCacheTask) Handel(message string) error {
	time.Sleep(time.Millisecond * 100)
	p := queue.DelCachePayload{}
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
	return 100
}

func (t *DelCacheTask) GetRetryCount() int {
	return 3
}

func (t *DelCacheTask) GetRateLimit() int {
	return 200
}
