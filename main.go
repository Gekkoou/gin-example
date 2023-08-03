package main

import (
	"fmt"
	"gin-example/initialize"
)

func main() {

	initialize.Initialize(false)

	// 初始化翻译器
	if err := initialize.InitTrans("zh"); err != nil {
		fmt.Printf("初始化翻译器错误, err = %s", err.Error())
		return
	}
	r := initialize.Routes()
	r.Run("127.0.0.1:8888")
}
