package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

/**
 * 下载驱动
 * go install github.com/go-sql-driver/mysql
 */

var db *sql.DB

func initDriverDB() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/sql_test?charset=utf8"
	// 这里只会连接 不会校验用户名和密码
	if db, err = sql.Open("mysql", dsn); err != nil {
		return
	}
	// ping 校验用户名和密码
	if err = db.Ping(); err != nil {
		return
	}

	// 设置数据库连接池最大连接数
	db.SetMaxOpenConns(10)
	// 设置数据库连接池最大空闲连接数
	db.SetMaxIdleConns(5)

	return
}

/**
* 创建 sql_test 库下 user 表
* use sql_test
CREATE TABLE `user` (
 `id` int(11) NOT NULL AUTO_INCREMENT,
 `name` varchar(255) DEFAULT '',
 `age` int(11) DEFAULT '0',
 PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
*/

type User struct {
	ID   int
	Name string
	Age  int
}

func DriverDo() {
	if err := initDriverDB(); err != nil {
		fmt.Println("init db failed, err:", err)
		return
	}

	fmt.Println("mysql connect success")

	//if err := queryOne(1); err != nil {
	//	fmt.Println("queryOne failed, err:", err)
	//	return
	//}
	//if err := queryMany(0); err != nil {
	//	fmt.Println("queryMany failed, err:", err)
	//	return
	//}
	//if err := insert("anna", 30); err != nil {
	//	fmt.Println("insert failed, err:", err)
	//	return
	//}
	//if err := update(1, 25); err != nil {
	//	fmt.Println("update failed, err:", err)
	//	return
	//}
	//if err := delete(3); err != nil {
	//	fmt.Println("delete failed, err:", err)
	//	return
	//}
	//if err := prepareInsert(); err != nil {
	//	fmt.Println("prepareInsert failed, err:", err)
	//	return
	//}
	if err := transaction(); err != nil {
		fmt.Println("transaction failed, err:", err)
		return
	}
}

// 查询单条
func queryOne(id int) error {
	var u User
	sqlStr := "SELECT id,name,age FROM user WHERE id=?"
	row := db.QueryRow(sqlStr, id)
	if err := row.Scan(&u.ID, &u.Name, &u.Age); err != nil {
		return err
	}
	fmt.Println(u)

	return nil
}

// 查询多条
func queryMany(id int) error {
	var us []User
	sqlStr := "SELECT id,name,age FROM user WHERE id>?"
	rows, err := db.Query(sqlStr, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		if err = rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
			return err
		}
		us = append(us, u)
	}

	fmt.Println(us)

	return nil
}

// 插入
func insert(name string, age int) error {
	sqlStr := "INSERT INTO user(name,age) VALUES(?,?)"
	ret, err := db.Exec(sqlStr, name, age)
	if err != nil {
		return err
	}

	lastInsertId, err := ret.LastInsertId()
	if err != nil {
		return err
	}
	fmt.Println(lastInsertId)

	return nil
}

// 更新
func update(id, age int) error {
	sqlStr := "UPDATE user SET age=? WHERE id=?"
	ret, err := db.Exec(sqlStr, age, id)
	if err != nil {
		return err
	}

	rowsAffected, err := ret.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println(rowsAffected)

	return nil
}

// 删除
func delete(id int) error {
	sqlStr := "DELETE FROM user WHERE id=?"
	ret, err := db.Exec(sqlStr, id)
	if err != nil {
		return err
	}

	rowsAffected, err := ret.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println(rowsAffected)

	return nil
}

// 预处理插入
func prepareInsert() error {
	sqlStr := "INSERT INTO user(name,age) VALUES(?,?)"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	m := map[string]int{
		"tom":   20,
		"jerry": 25,
	}
	for name, age := range m {
		stmt.Exec(name, age)
	}

	return nil
}

// 事务
func transaction() error {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("begin transaction failed, err:", err)
		return err
	}

	sqlStr := "INSERT INTO user(name,age) VALUES(?,?)"
	stmt, err := tx.Prepare(sqlStr)
	if err != nil {
		fmt.Println("prepare failed, err:", err)
		return err
	}
	defer stmt.Close()

	m := map[string]int{
		"john":  23,
		"marry": 24,
	}
	for name, age := range m {
		stmt.Exec(name, age)
	}

	if err = tx.Commit(); err != nil {
		fmt.Println("commit failed, err:", err)
		return err
	}

	// 如果出错了 回滚
	//tx.Rollback()

	fmt.Println("transaction success")
	return nil
}
