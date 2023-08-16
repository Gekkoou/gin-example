package initialize

import (
	"gin-example/global"
	_ "gin-example/queue"
	"gin-example/queue/core"
)

func NewQueue() (*core.Queue, error) {
	// 初始化队列
	core.QueueApp.Logger = core.LoggerFunc(NewQueueLogger)
	core.QueueApp.ErrorLogger = core.LoggerFunc(NewQueueErrorLogger)
	err := core.QueueApp.NewJob(global.Config.Queue)
	return core.QueueApp, err
}

func NewQueueLogger(msg string, args ...interface{}) {
	global.Log.Info(msg)
}

func NewQueueErrorLogger(msg string, args ...interface{}) {
	global.Log.Error(msg)
}
