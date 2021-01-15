/********************
@File        : record
@Time        : 2020-12-15 17:24:21
@Author      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : 狼人杀战绩数据
**********************/

package Mydb

import (
	"log"
	"time"
)

type Record struct {
	Id       int64
	User     int       `xorm:"User_id`
	GameTime time.Time `xorm:"text"`
	Identity int       `xorm:"integer"`
	GameMode string    `xorm:"varchar(255)"`
	RunAway  int       `xorm:"integer"`
	Result   string    `xorm:"varchar(255)"`
}

// 获取全部战绩
func (r Record) GetUser(a ...interface{}) []Record {
	r, ok := a[0].(Record)
	records := make([]Record, 0)
	if ok != false {
		err := orm.Find(&records, r)
		if err != nil {
			log.Panic(err)
		}
	}
	return records
}

// 插入单个战绩
func (r Record) Insert(a ...interface{}) bool {
	_, err := orm.InsertOne(a[0])
	if err != nil {
		log.Panic(err)
	}
	return true
}

// 战绩不允许修改 不允许隐藏
func (r Record) Update(a ...interface{}) bool {
	return false
}
