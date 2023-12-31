package config

type Mysql struct {
	Path         string `mapstructure:"path" json:"path"`                     // 服务器地址:端口
	Port         string `mapstructure:"port" json:"port"`                     // :端口
	Config       string `mapstructure:"config" json:"config"`                 // 高级配置
	Dbname       string `mapstructure:"db-name" json:"db-name"`               // 数据库名
	Username     string `mapstructure:"username" json:"username"`             // 数据库用户名
	Password     string `mapstructure:"password" json:"password"`             // 数据库密码
	Prefix       string `mapstructure:"prefix" json:"prefix"`                 // 全局表前缀，单独定义TableName则不生效
	Singular     bool   `mapstructure:"singular" json:"singular"`             // 是否开启全局禁用复数，true表示开启
	Engine       string `mapstructure:"engine" json:"engine" yaml:"engine"`   // 数据库引擎，默认InnoDB
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"max-open-conns"` // 打开到数据库的最大连接数
	LogMode      string `mapstructure:"log-mode" json:"log-mode"`             // 是否开启Gorm全局日志
	ZapLogLevel  string `mapstructure:"zap-log-level" json:"zap-log-level"`   // zap-log 的级别
}

func (m *Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + m.Port + ")/" + m.Dbname + "?" + m.Config
}
