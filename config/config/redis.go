package config

type Redis struct {
	DB       int    `mapstructure:"db" json:"db"`             // redis的哪个数据库
	Addr     string `mapstructure:"addr" json:"addr"`         // 服务器地址:端口
	Password string `mapstructure:"password" json:"password"` // 密码
}
