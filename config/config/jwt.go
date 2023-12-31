package config

type JWT struct {
	SigningKey  string `mapstructure:"signing-key" json:"signing-key"`   // 签名密钥
	ExpiresTime string `mapstructure:"expires-time" json:"expires-time"` // 过期时间
	Issuer      string `mapstructure:"issuer" json:"issuer"`             // 签发者
	TokenKey    string `mapstructure:"token-key" json:"token-key"`
}
