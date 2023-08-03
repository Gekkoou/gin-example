package config

type Serve struct {
	Ip   string `mapstructure:"ip" json:"ip" yaml:"ip"`       // Host
	Addr string `mapstructure:"addr" json:"addr" yaml:"addr"` // 服务器地址:端口
}
