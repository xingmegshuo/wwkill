/********************
@File        : buddy
@Time        : 2020-12-15 17:21:05
@Author      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : 狼人杀好友关系
**********************/

package Mydb

import "log"

type Buddy struct {
	Id   int64
	User int `xorm:"foregin key(user) references user(userid)"`
	// Own_user int `xorm:"foregin key(user) references user(userid)"`
	Agree  int    `xorm:"integer"`
	Buddys string `xorm:"text"`
	Del    int    `xorm:"integer"`
}

// 获取全部好友
func (b Buddy) GetUser(a ...interface{}) []Buddy {
	b, ok := a[0].(Buddy)
	// log.Println(b)
	buddys := make([]Buddy, 0)
	if ok != false {
		err := orm.Find(&buddys, b)
		if err != nil {
			log.Panic(err)
		}
	}
	// log.Println(buddys)
	return buddys
}

// 获取单个好友
func (b Buddy) GetBuddy(a ...interface{}) (Buddy, bool) {
	u, ok := a[0].(Buddy)
	if ok != false {
		has, _ := orm.Get(&u)
		return u, has
	} else {
		return Buddy{}, false
	}
}

// 插入单个好友
func (b Buddy) Insert(a ...interface{}) bool {
	// log.Println(a[0])
	_, err := orm.InsertOne(a[0])
	if err != nil {
		log.Panic(err)
	}
	return true
}

// 也就是删除单个好友 好友删除为双向，调用此方法时应该互相删除
func (b Buddy) Update(a ...interface{}) bool {
	b, ok := a[0].(Buddy)
	if ok != false {
		_, err := orm.Id(b.Id).Update(b)
		if err != nil {
			log.Panic(err)
		}
	}
	return true
}
