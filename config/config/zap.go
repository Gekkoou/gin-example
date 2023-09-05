package config

import (
	"go.uber.org/zap/zapcore"
	"strings"
)

type Zap struct {
	Level         string `mapstructure:"level" json:"level"`                   // 级别
	Prefix        string `mapstructure:"prefix" json:"prefix"`                 // 日志前缀
	Format        string `mapstructure:"format" json:"format"`                 // 输出
	Directory     string `mapstructure:"directory" json:"directory"`           // 日志文件夹
	EncodeLevel   string `mapstructure:"encode-level" json:"encode-level"`     // 编码级
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktrace-key"` // 栈名

	MaxAge       int  `mapstructure:"max-age" json:"max-age"`               // 日志留存时间
	ShowLine     bool `mapstructure:"show-line" json:"show-line"`           // 显示行
	LogInConsole bool `mapstructure:"log-in-console" json:"log-in-console"` // 输出控制台
}

func (z *Zap) ZapEncodeLevel() zapcore.LevelEncoder {
	switch {
	case z.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		return zapcore.LowercaseLevelEncoder
	case z.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		return zapcore.LowercaseColorLevelEncoder
	case z.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		return zapcore.CapitalLevelEncoder
	case z.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

func (z *Zap) TransportLevel() zapcore.Level {
	z.Level = strings.ToLower(z.Level)
	switch z.Level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.WarnLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}
