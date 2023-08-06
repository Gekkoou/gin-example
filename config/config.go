package config

import "gin-example/config/config"

type Config struct {
	Mysql config.Mysql
	Zap   config.Zap
	Redis config.Redis
	Serve config.Serve
	JWT   config.JWT
}
