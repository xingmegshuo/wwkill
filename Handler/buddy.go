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
	"math/rand"
	"strconv"
	"strings"
	"time"
	"wwKill/Mydb"
)

// 用户好友
func GetBuddy(mes []byte) string {
	err := json.Unmarshal(mes, &user)
	if err != nil {
		log.Println("数据问题:", err.Error())
		return ToMes("error", "获取好友失败,数据无法解析")
	}
	User := Mydb.User{
		OpenID: user.OpenID,
	}
	thisUser, has := ctrlUser.GetUser(User)
	if has {
		buddy := Mydb.Buddy{
			User:  int(thisUser.Id),
			Agree: 1,
			Del:   0,
		}
		backs := ctrlBuddy.GetUser(buddy)
		return BuddyToString("ok", backs, "获取好友列表成功", 0)
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
	User := Mydb.User{
		OpenID: user.OpenID,
	}
	thisUser, has := ctrlUser.GetUser(User)
	if has {
		back := Mydb.Buddy{
			Buddys:  strconv.FormatInt(thisUser.Id, 10),
			Agree:  0,
			Del:    0,
		}
		// log.Println(back)
		backs := ctrlBuddy.GetUser(back)
		return BuddyToString("ok", backs, "获取好友申请成功", 1)
	} else {
		return ToMes("error", "获取好友申请失败，找不到用户")
	}
}

// 同意好友申请
func AgreeBuddy(mes []byte) string {
	err := json.Unmarshal(mes, &buddy)
	if err != nil {
		log.Println("数据问题:", err.Error())
		return ToMes("error", "同意好友失败,数据无法解析")
	}
	B,_ := ctrlBuddy.GetBuddy(buddy)
	B.Agree = 1
	ctrlBuddy.Update(B)
	B, has := ctrlBuddy.GetBuddy(B)
	if has {
		UserId,_ :=  strconv.Atoi(B.Buddys)
		Another := Mydb.Buddy{
			User:  UserId,
			Buddys: strconv.Itoa(B.User),
			Agree:  1,
		}
		ctrlBuddy.Insert(Another)
	}

	return ToMes("ok", "同意好友成功")
}

// 删除好友
func DeleteBuddy(mes []byte) string {
	err := json.Unmarshal(mes, &buddy)
	if err != nil {
		log.Println("数据问题:", err.Error())
		return ToMes("error", "删除好友失败,数据无法解析")
	}
	Newbuddy := Mydb.Buddy{
		Id:  buddy.Id,
		Del: 1,
	}
	ctrlBuddy.Update(Newbuddy)
	B, has := ctrlBuddy.GetBuddy(buddy)
	if has {
		userId, _ := strconv.ParseInt(B.Buddys, 10, 64)
		Another := Mydb.Buddy{
			User:   int(userId),
			Buddys: strconv.Itoa(B.User),
		}
		a, _ := ctrlBuddy.GetBuddy(Another)
		this := Mydb.Buddy{
			Id:  a.Id,
			Del: 1,
		}
		ctrlBuddy.Update(this)
	}
	return ToMes("ok", "删除好友成功")

}

// 获取推荐好友
func RecomBuddy(mes []byte) string {
	err := json.Unmarshal(mes, &user)
	if err != nil {
		log.Println("数据问题:", err.Error())
		return ToMes("error", "同意好友失败,数据无法解析")
	}
	U, _ := ctrlUser.GetUser(user)
	SearchUser := Mydb.User{
		Id: 0,
	}
	users := ctrlUser.GetUsers(SearchUser)
	str := "{'status':'ok','mes':'获取推荐好友','data':["
	rand.Seed(time.Now().Unix())
	l := 0
	if len(users) > 3 {
		l = 3
	} else {
		l = len(users) - 1
	}
	for i := 0; i < l; i++ {
		for {
			j := rand.Intn(len(users))
			if users[j].Id != U.Id {
				if i == l-1 {
					str = str + "{'openID':'" + users[j].OpenID + "','nickName':'" + users[j].NickName + "','avatarUrl':'" + users[j].AvatarURL + "','level':'" + strconv.Itoa(users[j].Level) + "','id':'" + strconv.Itoa(int(users[j].Id)) + "'}"
				} else {
					str = str + "{'openID':'" + users[j].OpenID + "','nickName':'" + users[j].NickName + "','avatarUrl':'" + users[j].AvatarURL + "','level':'" + strconv.Itoa(users[j].Level) + "','id':'" + strconv.Itoa(int(users[j].Id)) + "'},"
				}
				break
			} else {
				j = rand.Intn(len(users))
				continue
			}
		}
	}
	str = str + "]}"
	str = strings.Replace(str, "'", "\"", -1)
	return str
}

// 添加好友
func AddBuddy(mes []byte) string {
	err := json.Unmarshal(mes, &buddy)
	if err != nil {
		log.Println("数据问题:", err.Error())
		return ToMes("error", "发送好友申请失败,数据无法解析")
	}
	thisBuddy := Mydb.Buddy{
		User:   int(buddy.User),
		Buddys: buddy.Buddys,
		Agree:  0,
	}
	B, has := ctrlBuddy.GetBuddy(buddy)
	if has {
		if B.Agree == 0 {
			return ToMes("error", "您已经向该用户发送过好友申请")
		} else {
			return ToMes("error", "该好友已经是您的好友")
		}
	} else {
		ctrlBuddy.Insert(thisBuddy)
		return ToMes("ok", "发送好友申请成功")
	}

}

// buddy to str
func BuddyToString(status string, back []Mydb.Buddy, mes string, statu int) string {
	str := "{'status':'" + status + "','mes':'" + mes + "','data':["
	for l, item := range back {
		// to do : 解析好友列表
		if len(item.Buddys) > 0 && item.Del == 0 {

			if statu == 0 {
				UserId, _ := strconv.ParseInt(item.Buddys, 10, 64)
				// log.Println(userId)
				user = Mydb.User{
					Id: UserId,
				}
			} else {
				userId := item.User
				// log.Println(userId)
				user = Mydb.User{
					Id: int64(userId),
				}
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
