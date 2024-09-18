package initialize

import (
	"fmt"
	"gin-example/global"
	"gin-example/utils"
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
	// 初始化数据库
	err = InitDatabase()
	if err != nil {
		fmt.Printf("初始化数据库失败, err = %s", err.Error())
		return
	}
	// 初始化DB
	global.DB, err = GormMysql()
	if err != nil {
		fmt.Printf("Db初始化失败, err = %s", err.Error())
		return
	}
	// 迁移
	err = InitTable()
	if err != nil {
		fmt.Printf("迁移失败, err = %s", err.Error())
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

	go func() {
		for {

			time.Sleep(10 * time.Millisecond)
			go func() {
				utils.AsynQueue("del-cache", struct {
					Key string
				}{"del-cache-key-" + time.Now().Format(time.DateTime)})
			}()

			go func() {
				utils.AsynQueue("del-cache", struct {
					Key string
				}{"del-cache-key-" + time.Now().Format(time.DateTime)})
			}()

			go func() {
				utils.AsynQueue("del-cache", struct {
					Key string
				}{"del-cache-key-" + time.Now().Format(time.DateTime)})
			}()

			go func() {
				utils.AsynQueue("del-cache", struct {
					Key string
				}{"del-cache-key-" + time.Now().Format(time.DateTime)})
			}()

			go func() {
				utils.AsynQueue("del-cache", struct {
					Key string
				}{"del-cache-key-" + time.Now().Format(time.DateTime)})
			}()

			go func() {
				utils.AsynQueue("del-cache", struct {
					Key string
				}{"del-cache-key-" + time.Now().Format(time.DateTime)})
			}()

			go func() {
				utils.AsynQueue("test-job", struct {
					Key string
				}{"test-job" + time.Now().Format(time.DateTime)})
			}()

			go func() {
				utils.AsynQueue("test-job", struct {
					Key string
				}{"test-job" + time.Now().Format(time.DateTime)})
			}()

			go func() {
				utils.AsynQueue("test-job", struct {
					Key string
				}{"test-job" + time.Now().Format(time.DateTime)})
			}()

		}
	}()
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
