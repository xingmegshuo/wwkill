/***************************
@File        : game.go
@Time        : 2021/01/19 14:40:30
@AUTHOR      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : 游戏逻辑
****************************/

package Handler

import (
	"encoding/json"
	"log"
)

type GameType struct {
	GameMode string
}

type Room struct {
}

var Game GameType

// 开始游戏
func GameStart(mes []byte) string {
	log.Println("开始匹配")
	err := json.Unmarshal(mes, &Game)
	if err != nil {
		log.Println("数据问题:", err.Error())
		return ToMes("error", "数据无法解析")
	}
	switch Game.GameMode {
	case "普通场":
		log.Println("hhh")
	case "高手场":
		log.Println("hhh")

	default:
		log.Println("hhhh")

	}
	return "ok"
}
