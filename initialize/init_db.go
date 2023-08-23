package initialize

import (
	"database/sql"
	"fmt"
	"gin-example/global"
)

func InitDatabase() error {
	return initDatabase(getDsn(), "mysql", getCreateSql(global.Config.Mysql.Dbname), getCreateSql("dtm_barrier"), initDtmBarrier())
}

func getDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/", global.Config.Mysql.Username, global.Config.Mysql.Password, global.Config.Mysql.Path, global.Config.Mysql.Port)
}

func getCreateSql(name string) string {
	return fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", name)
}

func initDatabase(dsn string, driver string, sqls ...string) error {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(db)
	if err = db.Ping(); err != nil {
		return err
	}

	for _, sql := range sqls {
		if _, err = db.Exec(sql); err != nil {
			return err
		}
	}
	return nil
}

func initDtmBarrier() string {
	return `
CREATE TABLE IF NOT EXISTS dtm_barrier.barrier (
  id bigint(22) NOT NULL AUTO_INCREMENT,
  trans_type varchar(45) DEFAULT '',
  gid varchar(128) DEFAULT '',
  branch_id varchar(128) DEFAULT '',
  op varchar(45) DEFAULT '',
  barrier_id varchar(45) DEFAULT '',
  reason varchar(45) DEFAULT '' COMMENT 'the branch type who insert this record',
  create_time datetime DEFAULT CURRENT_TIMESTAMP,
  update_time datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY gid (gid,branch_id,op,barrier_id),
  KEY create_time (create_time),
  KEY update_time (update_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`
}
