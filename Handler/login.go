/***************************
@File        : login.go
@Time        : 2020/12/21 14:33:25
@AUTHOR      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : 解析登录
****************************/

package Handler

import (
	"encoding/json"
	"log"
	"wwKill/Mydb"
)

// 用户登录函数处理
func Login(mes []byte) string {
	ctrlUser := Mydb.NewUserCtrl()
	var user Mydb.User
	err := json.Unmarshal(mes, &user)
	if err != nil {
		log.Println("数据问题:", err.Error())
		return "登录操作失败"
	}

	log.Println(user.OpenID, user.NickName, ctrlUser)
	return "登录操作成功"
}
