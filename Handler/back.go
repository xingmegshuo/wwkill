/***************************
@File        : back.go
@Time        : 2021/01/12 11:43:01
@AUTHOR      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : 用户背包
****************************/

package Handler

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"wwKill/Mydb"
)

// 用户购买服装
func AddBack(mes []byte) string {
	err := json.Unmarshal(mes, &backpack)
	if err != nil {
		log.Println("数据问题:", err.Error())
		return ToMes("error", "获取背包失败,数据无法解析")
	}
	ctrlBack.Insert(backpack)
	return ToMes("ok", "添加背包成功")
}

// 用户背包
func GetBack(mes []byte) string {
	err := json.Unmarshal(mes, &user)
	if err != nil {
		log.Println("数据问题:", err.Error())
		return ToMes("error", "获取背包失败,数据无法解析")
	}
	User := Mydb.User{
		OpenID: user.OpenID,
	}
	thisUser, has := ctrlUser.GetUser(User)
	if has {
		back := Mydb.Backpack{
			User: int(thisUser.Id),
		}
		backs := ctrlBack.GetUser(back)
		return BackToString("ok", backs, "获取背包成功")
	} else {
		return ToMes("error", "获取背包失败，找不到用户")
	}

}

// back to str
func BackToString(status string, back []Mydb.Backpack, mes string) string {
	str := "{'status':'" + status + "','mes':'" + mes + "','data':["
	for l, item := range back {
		if l == len(back)-1 {
			str = str + "{'name':'" + item.Name + "','property':'" + strconv.Itoa(item.Property) + "','stilTime':'" + item.StilTime + "'}"
		} else {
			str = str + "{'name':'" + item.Name + "','property':'" + strconv.Itoa(item.Property) + "','stilTime':'" + item.StilTime + "'},"
		}
	}
	str = str + "]}"
	str = strings.Replace(str, "'", "\"", -1)
	return str
}
