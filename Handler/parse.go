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

	default:
		log.Println("无效key")
		mes := "没有对应数据处理!"
		return mes
	}
}
