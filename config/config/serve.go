package config

type Serve struct {
	Ip   string `mapstructure:"ip" json:"ip"`     // Host
	Addr string `mapstructure:"addr" json:"addr"` // 服务器地址:端口
}
