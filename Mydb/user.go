/***************************
@File        : user.go
@Time        : 2020/12/10 15:13:40
@AUTHOR      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : orm user
****************************/

package Mydb

import (
	"fmt"
	"log"
)

type User struct {
	Id        int64
	OpenID    string `xorm:"varchar(255)"`
	NickName  string `xorm:"varchar(255)"`
	AvatarURL string `xorm:"varchar(255)"`
	Level     int    `xorm:default 0`
	Money     int    `xorm:default 300`
	Orther    string `xorm:text`
}

// 根据条件返回用户
func (u User) GetUser(a ...interface{}) (User, bool) {
	u, ok := a[0].(User)
	if ok != false {
		has, _ := orm.Get(&u)
		fmt.Println(u, has)
		return u, has
	} else {
		return User{}, false
	}
}

// 插入单个用户
func (u User) Insert(a ...interface{}) bool {
	_, err := orm.InsertOne(a[0])
	if err != nil {
		log.Panic(err)
	}
	return true
}

// 修改
func (u User) Update(a ...interface{}) bool {
	u, ok := a[0].(User)
	if ok != false {
		_, err := orm.Id(u.Id).Update(u)
		if err != nil {
			log.Panic(err)
		}
	}
	return true
}
