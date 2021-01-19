/***************************
@File        : parse.go
@Time        : 2020/12/21 15:56:57
@AUTHOR      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : 解析数据根据不同key来进行区分
****************************/

package Handler

import (
	"encoding/json"
	"log"

	"golang.org/x/net/websocket"
)

// 数据格式
type Data struct {
	Name   string
	Values string
}

var client_user = make(map[*websocket.Conn]string)

// var client_buddy = make(map[*websocket.Conn]string)

// 解析key
func ParseData(con string, ws *websocket.Conn) string {
	log.Println("开始解析", con)
	var data Data
	oldData := []byte(con)
	err := json.Unmarshal(oldData, &data)
	if err != nil {
		log.Println(err)
		return "解析数据失败,请查看数据格式"
	}
	log.Printf("类型%T", data.Values)
	info := []byte(data.Values)
	// log.Println(info)
	// log.Println(data.Values, data)
	switch data.Name {
	case "login":
		log.Println("登录操作:")
		mes := Login(info, ws)
		return mes
	case "upgrade":
		log.Println("账号升级")
		mes := Upgrade(info)
		return mes
	case "back":
		log.Println("获取背包")
		mes := GetBack(info)
		return mes
	case "addback":
		log.Println("购买商品，增加背包")
		mes := AddBack(info)
		return mes
	case "record":
		log.Println("获取最近战绩")
		mes := GetRecord(info)
		return mes
	case "recordRate":
		log.Println("获取全部战斗")
		mes := GetRecordAll(info)
		return mes
	case "buddy":
		log.Println("获取好友列表")
		mes := GetBuddy(info)
		return mes
	case "newbuddy":
		log.Println("获取好友申请")
		mes := GetNewBuddy(info)
		return mes
	case "agreebuddy":
		log.Println("同意好友申请")
		mes := AgreeBuddy(info)
		return mes
	case "rcombuddy":
		log.Println("获取好友推荐")
		mes := RecomBuddy(info)
		return mes
	case "addbuddy":
		log.Println("添加好友申请")
		mes := AddBuddy(info)
		return mes
	case "delbuddy":
		log.Println("删除好友")
		mes := DeleteBuddy(info)
		return mes
	case "chat":
		log.Println("好友聊天")
		mes := Chat(info)
		return mes

	default:
		log.Println("无效key")
		mes := "没有对应数据处理!"
		return mes
	}
}

// 关闭连接时退出用户
func CloseUser(ws *websocket.Conn) {
	delete(client_user, ws)
}
