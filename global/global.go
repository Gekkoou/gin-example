package global

import (
	"gin-example/config"
	"github.com/go-playground/universal-translator"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	Redis  *redis.Client
	Trans  ut.Translator
	Config config.Config
	VP     *viper.Viper
	Log    *zap.Logger
)
