/***************************
@File        : upgrade.go
@Time        : 2021/01/11 13:36:51
@AUTHOR      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : 用户升级增加金币操作
****************************/

package Handler

import (
	"encoding/json"
	"log"
)

// 用户升级操作
func Upgrade(mes []byte) string {
	// ctrlUser := Mydb.NewUserCtrl()
	// var user Mydb.User
	err := json.Unmarshal(mes, &user)
	if err != nil {
		log.Println("数据问题:", err.Error())
		return ToMes("error", "升级失败,数据无法解析")
	}
	thisUser, has := ctrlUser.GetUser(user)
	if has {
		thisUser.Level = thisUser.Level + 1
		if user.Money != thisUser.Money && user.Money != 0 {
			thisUser.Money = user.Money
		}
		if len(user.NickName) > 0 {
			thisUser.NickName = user.NickName
		}
		if len(user.AvatarURL) > 0 {
			thisUser.AvatarURL = user.AvatarURL
		}
		num := ctrlUser.Update(thisUser)
		log.Println(num)
	} else {
		return ToMes("error", "升级失败，找不到用户")
	}
	return UserToString("ok", thisUser, "升级成功")
}
