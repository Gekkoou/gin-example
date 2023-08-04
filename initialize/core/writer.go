package core

import (
	"fmt"
	"gin-example/global"
)

type Writer struct {
}

// Printf 格式化打印日志
func (w Writer) Printf(message string, data ...interface{}) {
	logLevel := global.Config.Mysql.ZapLogLevel
	switch logLevel {
	case "error", "Error", "ERROR":
		global.Log.Error(fmt.Sprintf(message+"\n", data...))
	default:
		global.Log.Info(fmt.Sprintf(message+"\n", data...))
	}
}
