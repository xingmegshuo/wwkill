/***************************
@File        : table.go
@Time        : 2020/12/10 15:53:52
@AUTHOR      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        :  创建数据库
****************************/

package Mydb

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// 生成创建数据表的sql
func table(n int) (s string) {
	switch {
	// 创建 游戏用户数据表
	case n == 0:
		s := `create table if not exists "users" (id integer 
			primary key, openId varchar(255),level integer default 0,
			 money integer default 300, orther text );`
		return s
		// 个人背包，在商店买的东西
	case n == 1:
		s := `create table if not exists "store" (id integer 
			primary key, name varchar(255),price integer,
			property integer default 0, num interger, user interger, 
			foreign key(user) references users(usersid));`
		return s
		// 好友关系
	case n == 2:
		s := `create table if not exists "buddy" (id integer primary key,
		user integer, own_user integer , argree integer,
		foreign key(user) references users(userssid),
		foreign key(own_user) references users(usersid));`
		return s
		// 战绩信息
	case n == 3:
		s := `create table if not exists "record" (
			id integer primary key, user integer, gameTime text,
			identity varchar(255), gameMode varchar(255),count integer,
			winCount integer, runAway integer, winRate integer, result integer,
			foreign key(user) references users(usersid)
		);`
		return s
	default:
		s := `select * from sqlite_master where type="table";`
		return s
	}
}

// 初始化数据库
func Db(dbName string) *sql.DB {
	db, err := sql.Open("sqlite3", dbName+".db")
	if err != nil {
		fmt.Printf("错误:%s", err)
	}
	for i := 0; i < 4; i++ {
		s := table(i)
		_, err = db.Exec(s)
		if err != nil {
			fmt.Println(err)
		}
	}
	return db
}
