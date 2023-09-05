package initialize

import (
	"fmt"
	"gin-example/global"
	"gin-example/initialize/core"
	"gin-example/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// Zap 获取 zap.Logger
func Zap() (logger *zap.Logger, err error) {
	if ok, _ := utils.PathExists(global.Config.Zap.Directory); !ok { // 判断是否有Directory文件夹
		fmt.Printf("create %v directory\n", global.Config.Zap.Directory)
		_ = os.Mkdir(global.Config.Zap.Directory, os.ModePerm)
	}
	cores := core.Zap.GetZapCores()
	logger = zap.New(zapcore.NewTee(cores...))
	if global.Config.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return
}
