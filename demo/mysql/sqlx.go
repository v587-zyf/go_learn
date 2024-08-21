package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

/**
 * go install github.com/jmoiron/sqlx
 */

var sqlDB *sqlx.DB

func initSqlxDB() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/sql_test?charset=utf8"

	if sqlDB, err = sqlx.Connect("mysql", dsn); err != nil {
		return
	}

	// 设置数据库连接池最大连接数
	sqlDB.SetMaxOpenConns(10)
	// 设置数据库连接池最大空闲连接数
	sqlDB.SetMaxIdleConns(5)

	return
}

func SqlxDo() {
	if err := initSqlxDB(); err != nil {
		fmt.Println("sqlx init err:", err)
		return
	}

	sqlStr := "SELECT id,name,age FROM user WHERE id = ?"
	var u User
	if err := sqlDB.Get(&u, sqlStr, 1); err != nil {
		fmt.Println("sqlx Get err:", err)
		return
	}
	fmt.Println(u)

	sqlStr = "SELECT id,name,age FROM user"
	var uList []User
	if err := sqlDB.Select(&uList, sqlStr); err != nil {
		fmt.Println("sqlx Select err:", err)
		return
	}
	fmt.Println(uList)
}
