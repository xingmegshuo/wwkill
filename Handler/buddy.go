/***************************
@File        : buddy.go
@Time        : 2021/01/12 16:28:06
@AUTHOR      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : 处理好友关系
****************************/

package Handler

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"wwKill/Mydb"
)

// 用户好友
func GetBuddy(mes []byte) string {
	err := json.Unmarshal(mes, &user)
	if err != nil {
		log.Println("数据问题:", err.Error())
		return ToMes("error", "获取好友失败,数据无法解析")
	}
	thisUser, has := ctrlUser.GetUser(user)
	if has {
		back := Mydb.Buddy{
			User:  int(thisUser.Id),
			Agree: 1,
		}
		backs := ctrlBack.GetUser(back)
		return BackToString("ok", backs, "获取好友列表成功")
	} else {
		return ToMes("error", "获取好友失败，找不到用户")
	}
}

// 获取好友申请
func GetNewBuddy(mes []byte) string {
	err := json.Unmarshal(mes, &user)
	if err != nil {
		log.Println("数据问题:", err.Error())
		return ToMes("error", "获取好友失败,数据无法解析")
	}
	thisUser, has := ctrlUser.GetUser(user)
	if has {
		back := Mydb.Buddy{
			User:  int(thisUser.Id),
			Agree: 0,
		}
		backs := ctrlBack.GetUser(back)
		return BackToString("ok", backs, "获取好友申请成功")
	} else {
		return ToMes("error", "获取好友失败，找不到用户")
	}
}

// 同意好友申请
func AgreeBuddy(mes []byte) string {
	err := json.Unmarshal(mes, &buddy)
	if err != nil {
		log.Println("数据问题:", err.Error())
		return ToMes("error", "同意好友失败,数据无法解析")
	}
	Newbuddy := Mydb.Buddy{
		Id:    buddy.Id,
		Agree: 1,
	}
	ctrlBuddy.Update(Newbuddy)
	return ToMes("ok", "同意好友成功")
}

// 获取推荐好友
// func RecomBuddy() string{

// }


// buddy to str
func BuddyToString(status string, back []Mydb.Buddy, mes string) string {
	str := "{'status':'" + status + "','mes':'" + mes + "','data':["
	for l, item := range back {
		// to do : 解析好友列表
		if len(item.Buddys) > 0 && item.Del == 0 && item.Agree == 1 {
			userId, _ := strconv.ParseInt(item.Buddys, 10, 64)
			user := Mydb.User{
				Id: userId,
			}
			itemId := strconv.FormatInt(item.Id, 10)
			User, _ := ctrlUser.GetUser(user)
			if l == len(back)-1 {
				str = str + "{'openID':'" + User.OpenID + "','nickName':'" + User.NickName + "','avatarUrl':'" + User.AvatarURL + "','level':'" + strconv.Itoa(User.Level) + "','Id':'" + itemId + "'}"
			} else {
				str = str + "{'openID':'" + User.OpenID + "','nickName':'" + User.NickName + "','avatarUrl':'" + User.AvatarURL + "','level':'" + strconv.Itoa(User.Level) + "','Id':'" + itemId + "'},"
			}
		}

	}
	str = str + "]}"
	str = strings.Replace(str, "'", "\"", -1)
	return str
}
