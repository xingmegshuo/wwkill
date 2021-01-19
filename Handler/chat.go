/***************************
@File        : chat.go
@Time        : 2021/01/19 10:36:21
@AUTHOR      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : 好友聊天
****************************/

package Handler

import (
	"encoding/json"
	"log"
	"wwKill/Mydb"

	"golang.org/x/net/websocket"
)

type ChatStruct struct {
	OpenId  string
	Message string
}

var chat ChatStruct

// 好友聊天
func Chat(info []byte) string {
	err := json.Unmarshal(info, &chat)
	if err != nil {
		log.Println("数据问题:", err.Error())
		return ToMes("error", "发送信息失败,数据解析问题")
	}
	User := Mydb.User{
		OpenID: chat.OpenId,
	}
	mes := chat.Message
	U, _ := ctrlUser.GetUser(User)
	var ws *websocket.Conn
	for k, v := range client_user {
		log.Println(v)
		if v == U.OpenID {
			ws = k
			break
		} else {
			ws = nil
			continue
		}
	}
	if ws == nil {
		return ToMes("error", "用户不在线")
	} else {
		err = websocket.Message.Send(ws, mes)
		if err != nil {
			return ToMes("error", "发送失败")
		} else {
			return ToMes("ok", "发送成功")
		}
	}
}
