package initialize

import (
	"errors"
	"gin-example/global"
	"gin-example/initialize/core"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GormMysql() (db *gorm.DB, err error) {
	m := global.Config.Mysql
	if m.Dbname == "" {
		return nil, errors.New("数据库不能为空")
	}
	mysqlConfig := mysql.Config{
		DSN:               m.Dsn(), // DSN data source name
		DefaultStringSize: 191,     // string 类型字段的默认长度
		//	SkipInitializeWithVersion: false,   // 根据版本自动配置
	}
	if db, err = gorm.Open(mysql.New(mysqlConfig), core.Gorm.Config(m.Prefix, m.Singular)); err != nil {
		return nil, err
	} else {
		// db.InstanceSet("gorm:table_options", "ENGINE="+m.Engine)
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
	}
	return
}
