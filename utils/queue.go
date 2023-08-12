package utils

import (
	"context"
	"fmt"
	"gin-example/global"
	"gin-example/queue/core"
	"github.com/bytedance/sonic"
)

func AsynQueue(task core.TaskInterFace, data interface{}) {
	msg, err := sonic.MarshalString(data)
	if err != nil {
		global.Log.Error(fmt.Sprintf("asynQueue marshal error: %s", err))
	}
	if err = global.Queue.Push(task, context.Background(), msg); err != nil {
		global.Log.Error(fmt.Sprintf("asynQueue push error: %s", err))
	}
}
