/********************
@File        : store
@Time        : 2020-12-15 17:10:35
@Author      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : 狼人杀个人背包数据库
**********************/

package Mydb

import (
	"log"
)

type Backpack struct {
	Id       int
	Name     string `xorm:"varchar(255)"`
	Property int    `xorm:"integer"`
	Num      int    `xorm:"integer"`
	StilTime string `xorm:"text"`
	User     int    `xorm:"foreign key(user) references user(userid)"`
	Del      int    `xorm:"integer"`
}

// 根据用户返回背包
func (b Backpack) GetUser(a ...interface{}) []Backpack {
	b, ok := a[0].(Backpack)
	backs := make([]Backpack, 0)
	if ok != false {
		err := orm.Find(&backs, b)
		if err != nil {
			log.Panic(err)
		}
	}
	return backs
}

// 插入单个物品
func (b Backpack) Insert(a ...interface{}) bool {
	_, err := orm.InsertOne(a[0])
	if err != nil {
		log.Panic(err)
	}
	return true
}

// 修改单个物品
func (b Backpack) Update(a ...interface{}) bool {
	// example : 服装到期, 小好评使用完
	b, ok := a[0].(Backpack)
	if ok != false {
		_, err := orm.Id(b.Id).Update(b)
		if err != nil {
			log.Panic(err)
		}
	}
	return true
}
