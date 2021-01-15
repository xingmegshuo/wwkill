/***************************
@File        : db_test.go
@Time        : 2020/12/10 15:24:54
@AUTHOR      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : 测试数据库
****************************/

package main

import (
	"fmt"
	"wwKill/Mydb"
)

func main() {

	// 背包管理器测试
	ctrl := Mydb.NewBuddyCtrl()
	// back := Mydb.Backpack{
	// 	Name: "可爱裤子",
	// 	User: 1,
	// }
	// ctrl.Insert(back)
	// back := Mydb.Backpack{
	// 	Id:       1,
	// 	StilTime: "2012-12-12 08：00：00",
	// }
	// ctrl.Update(back)
	buddy := Mydb.Buddy{
		User:  1,
		Agree: 1,
	}
	// ctrl.Insert(buddy)
	backs := ctrl.GetUser(buddy)
	fmt.Println(backs)

	// 用户管理器测试
	// ctrlUser := Mydb.NewUserCtrl()
	// user := Mydb.User{
	// 	OpenID:    "abcd",
	// 	NickName:  "abcd",
	// 	AvatarURL: "ddddd",
	// 	Orther:    "ahhh",
	// 	Money:     300,
	// }
	// // 插入一个
	// ctrlUser.Insert(user)
	// update := Mydb.User{Id: 1}
	// user, has := ctrlUser.GetUser(update)
	// fmt.Println(has)
	// if has {
	// 	fmt.Println("输出:", user)
	// }

}
