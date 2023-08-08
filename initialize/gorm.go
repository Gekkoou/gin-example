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
		DSN: m.Dsn(),
	}
	if db, err = gorm.Open(mysql.New(mysqlConfig), core.Gorm.Config(m.Prefix, m.Singular)); err != nil {
		return nil, err
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE="+m.Engine)
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
	}
	return
}
