package core

import (
	"fmt"
	"gin-example/global"
)

type Writer struct {
}

// Printf 格式化打印日志
func (w Writer) Printf(message string, data ...interface{}) {
	global.Log.Info(fmt.Sprintf(message+"\n", data...))
}
