package utils

import (
	"context"
	"fmt"
	"gin-example/global"
	"github.com/bytedance/sonic"
)

func AsynQueue(queue string, data any) {
	msg, err := sonic.MarshalString(data)
	if err != nil {
		global.Log.Error(fmt.Sprintf("asynQueue marshal error: %s", err))
	}
	if err = global.Queue.Push(queue, context.Background(), msg); err != nil {
		global.Log.Error(fmt.Sprintf("asynQueue push error: %s", err))
	}
}
