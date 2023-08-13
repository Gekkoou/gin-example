package initialize

import (
	"gin-example/global"
	_ "gin-example/queue"
	"gin-example/queue/core"
)

func NewQueue() (*core.Queue, error) {
	// 初始化队列
	core.QueueApp.Logger = global.Log
	err := core.QueueApp.NewJob(global.Config.Queue)
	return core.QueueApp, err
}
