package initialize

import (
	"context"
	"fmt"
	"gin-example/global"
	"gin-example/queue"
	"math/rand"
	"strconv"
	"time"
)

func Initialize(reload bool) {
	var err error
	// 配置热重载不需要执行
	if !reload {
		// VIPER
		global.VP, err = Viper()
		if err != nil {
			fmt.Printf("Viper 初始化失败, err = %s", err.Error())
			return
		}
	}
	// 日志
	global.Log, err = Zap()
	if err != nil {
		fmt.Printf("Log初始化失败, err = %s", err.Error())
		return
	}
	// 初始化DB
	global.DB, err = GormMysql()
	if err != nil {
		fmt.Printf("Db初始化失败, err = %s", err.Error())
		return
	}
	// 迁移
	err = RegisterTables()
	if err != nil {
		fmt.Printf("数据库迁移失败, err = %s", err.Error())
		return
	}
	// 初始化Redis
	global.Redis, err = Redis()
	if err != nil {
		fmt.Printf("Redis初始化失败, err = %s", err.Error())
		return
	}
	// 初始化翻译器
	err = InitTrans("zh")
	if err != nil {
		fmt.Printf("初始化翻译器错误, err = %s", err.Error())
		return
	}
	global.Queue, err = NewQueue()
	if err != nil {
		fmt.Printf("初始化队列, err = %s", err.Error())
		return
	}

	//  测试部分
	// 推送调用
	go func() {
		for true {
			if err := global.Queue.Push(queue.TestTaskApp, context.Background(), strconv.Itoa(rand.Intn(999999))); err != nil {
				fmt.Println(err)
				global.Log.Error(err.Error())
			}
			time.Sleep(time.Second)
		}
	}()
	// 测试部分结束

	// 注册路由
	r := Routes()
	if !reload {
		err = Run(r)
		if err != nil {
			fmt.Printf("启动失败, err = %s", err.Error())
			return
		}
	}
	return
}
