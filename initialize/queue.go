package initialize

import (
	"gin-example/global"
	"gin-example/queue/core"
)

func NewQueue() (*core.Queue, error) {
	// 初始化队列
	err := core.QueueApp.NewJob(global.Config.Queue)
	return core.QueueApp, err
}
