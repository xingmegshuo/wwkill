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
var client_palyer = make(map[*websocket.Conn]string)

// var client_buddy = make(map[*websocket.Conn]string)

// 解析key
func ParseData(con string, ws *websocket.Conn) {
	// log.Println("开始解析", con)
	var data Data
	oldData := []byte(con)
	err := json.Unmarshal(oldData, &data)
	if err != nil {
		log.Println(err)
		Send(ws, "解析数据失败,请查看数据格式")
	}
	// log.Printf("类型%T", data.Values)
	info := []byte(data.Values)
	switch data.Name {
	case "login":
		log.Println("登录操作:")
		mes := Login(info, ws)
		Send(ws, mes)
	case "upgrade":
		log.Println("账号升级")
		mes := Upgrade(info)
		Send(ws, mes)
	case "back":
		log.Println("获取背包")
		mes := GetBack(info)
		Send(ws, mes)
	case "addback":
		log.Println("购买商品，增加背包")
		mes := AddBack(info)
		Send(ws, mes)
	case "record":
		log.Println("获取最近战绩")
		mes := GetRecord(info)
		Send(ws, mes)
	case "recordRate":
		log.Println("获取全部战斗")
		mes := GetRecordAll(info)
		Send(ws, mes)
	case "buddy":
		log.Println("获取好友列表")
		mes := GetBuddy(info)
		Send(ws, mes)
	case "newbuddy":
		log.Println("获取好友申请")
		mes := GetNewBuddy(info)
		Send(ws, mes)
	case "agreebuddy":
		log.Println("同意好友申请")
		mes := AgreeBuddy(info)
		Send(ws, mes)
	case "rcombuddy":
		log.Println("获取好友推荐")
		mes := RecomBuddy(info)
		Send(ws, mes)
	case "addbuddy":
		log.Println("添加好友申请")
		mes := AddBuddy(info)
		Send(ws, mes)
	case "delbuddy":
		log.Println("删除好友")
		mes := DeleteBuddy(info)
		Send(ws, mes)
	case "chat":
		log.Println("好友聊天")
		mes := Chat(info)
		Send(ws, mes)
	case "getUser":
		log.Println("获取游戏中信息")
		mes := GetUserMes(info)
		Send(ws, mes)
	case "game":
		log.Println("开始游戏")
		mes := GameStart(info, ws)
		Send(ws, mes)
	default:
		log.Println("游戏中")
		go RoomSocket(oldData)
	}
}

// 关闭连接时退出用户
func CloseUser(ws *websocket.Conn) {
	delete(client_user, ws)
}

// 数据返回
func Send(ws *websocket.Conn, mes string) {
	if err := websocket.Message.Send(ws, mes); err != nil {
		log.Println("客户端丢失", err.Error())
		CloseUser(ws)
	}
}
