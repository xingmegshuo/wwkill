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
)

// 数据格式
type Data struct {
	Name   string
	Values string
}

// 解析key
func ParseData(con string) string {
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
	log.Println(info)
	// log.Println(data.Values, data)

	switch data.Name {
	case "login":
		log.Println("登录操作:")
		mes := Login(info)
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
	default:
		log.Println("无效key")
		mes := "没有对应数据处理!"
		return mes
	}
}
