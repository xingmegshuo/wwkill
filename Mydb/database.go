/***************************
@File        : database.go
@Time        : 2020/12/10 15:09:43
@AUTHOR      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : sqllite数据库
****************************/
package Mydb

import (
	"log"
	"time"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

// 数据库操作接口
type DbCtrlooer interface {
	Insert(a ...interface{}) bool
	Update(a ...interface{}) bool
	GetUser(a ...interface{})
}

var orm *xorm.Engine

// orm 实例返回
func SetEngin(db string) *xorm.Engine {
	var err error
	// 此处不能：= 返回的指针不能使用
	orm, err = xorm.NewEngine("sqlite3", db)
	if err != nil {
		log.Panic(err)
	}
	// 同步数据表
	this_err := orm.Sync2(new(Buddy), new(Record), new(User), new(Backpack))
	if this_err != nil {
		log.Panic("Fail to sync database: %v\n", this_err)
	}
	orm.TZLocation = time.Local
	return orm
}

// 初始化
func Init() {
	SetEngin("./wwKill.db")
}

// 用户管理器
func NewUserCtrl() *User {
	Init()
	return &User{}
}

// 背包管理器
func NewBackCtrl() *Backpack {
	Init()
	return &Backpack{}
}

// 好友管理器
func NewBuddyCtrl() *Buddy {
	Init()
	return &Buddy{}
}

// 战绩管理
func NewRecordCtrl() *Record {
	Init()
	return &Record{}
}
