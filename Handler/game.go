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
	"strconv"
	"time"

	"golang.org/x/net/websocket"
)

type GameType struct {
	GameMode string
	New      int
}

type Player struct {
	Identity string
	OpenID   string
	Survive  int
	Ws       *websocket.Conn
	Ready    int
}

type Room struct {
	GameMode string
	People   int
	Public   int
	User     []Player
	Owner    string
	Ww       int //狼人
	Ci       int //平民
	God      int //预言家
	Wi       int //女巫
	Hu       int //猎人
}

type Mes struct {
	Room    string
	Message string
	User    string
}

var PlayRoom = make(map[string]Room)

var Game GameType

// 开始游戏
func GameStart(mes []byte, ws *websocket.Conn) string {
	log.Println("进入游戏逻辑")
	var room Room
	err := json.Unmarshal(mes, &Game)
	if err != nil {
		log.Println("数据问题:", err.Error())
	}
	log.Println("进入新手房间")
	if Game.New == 0 {
		room = SearchRoom(Game.GameMode)
	} else {
		room = NewRoom(Game.GameMode)
	}
	message := Init(ws, room)
	return ToMes("ok", message)
}

// // 保持游戏中的连接,游戏中的通信
// func GameSocket(ws *websocket.Conn, room Room) {
// 	for {

// 	}
// }

// 查找房间
func SearchRoom(GameType string) Room {
	var room Room
	if len(PlayRoom) == 0 {
		room = NewRoom(GameType)
	} else {
		for _, v := range PlayRoom {
			if len(v.User) < v.People && v.Public == 1 && v.GameMode == GameType {
				room = v
				break
			} else {
				room = Room{
					Owner: "null",
				}
			}
		}

	}
	if room.Owner != "null" {
		return room
	} else {
		room = NewRoom(GameType)
		return room
	}

}

// 新建房间
func NewRoom(GameType string) Room {
	var room Room
	switch GameType {
	case "新手场":
		room = Room{
			GameMode: "新手场",
			People:   6,
			Public:   1,
			Ww:       2,
			God:      1,
			Ci:       2,
			Hu:       1,
		}
	case "普通场":
		room = Room{
			GameMode: "普通场",
			People:   10,
			Public:   1,
		}
	case "高手场":
		room = Room{
			GameMode: "高手场",
			People:   10,
			Public:   1,
		}
	}
	return room
}

// 游戏逻辑
func GameSocket(room Room) {
	for {
		if len(room.User) == room.People {
			ServerSend(room, "法官:人员已到齐，请所有人准备!")
			continue
		}
		ready := Ready(room)
		if ready == 1 {
			ServerSend(room, "法官:所有人已准备，游戏5秒后开始!")
			time.Sleep(time.Second * 5)
			ServerSend(room, "游戏开始！")
		}
		over := Over(room)
		if over == 1 {
			ServerSend(room, "法官:游戏结束，----胜利!")
			break
		}
	}
}

// 加入房间
func Join(room Room, player Player) Room {
	room.User = append(room.User, player)
	ServerSend(room, "用户"+player.OpenID+"进入房间")
	return room
}

// 离开房间
// func Leave(room Room, player Player) {	err := json.Unmarshal(info, &chat)

// 	for l, item := range room.User {
// 		if item.OpenID == player.OpenID {
// 			room.User = append(room.User[:l], room.User[l+1:])
// 			break
// 		}
// 	}
// }

// 房间连接
func RoomSocket(conn []byte) {
	var mes Mes
	err := json.Unmarshal(conn, &mes)
	if err != nil {
		log.Println("连接断开")
		// todo 退出
	} else {
		for _, room := range PlayRoom {
			if room.Owner == mes.Room {
				if mes.Message == "准备" {
					room = Re(room, mes.User)
				}
				if mes.Message[:2] == "身份" {
					// to do 分配身份
					room = Iden(room, mes.User, mes.Message[2:])

				}
				GameSocket(room)
			}
		}
	}

}

// 初始化改变用户状态,由在线换为游戏中
func Init(ws *websocket.Conn, room Room) string {
	// openId 赋值给游戏中用户
	client_palyer[ws] = client_user[ws]
	// 把用户从在线用户移除
	delete(client_user, ws)
	// go GameSocket(ws, room)
	player := Player{
		OpenID: client_palyer[ws],
		Ws:     ws,
	}
	if len(room.User) == 0 {
		room.Owner = player.OpenID
		room = Join(room, player)
	} else {
		room = Join(room, player)
	}
	PlayRoom[strconv.Itoa(len(PlayRoom)+1)] = room
	return room.Owner
}

// 退出匹配或者退出房间或者退出游戏
func Close(ws *websocket.Conn) {
	client_user[ws] = client_palyer[ws]
	delete(client_palyer, ws)
}

// 发送服务器公告信息
func ServerSend(room Room, str string) {
	for _, U := range room.User {
		Send(U.Ws, str)
	}
}

// 判断是否全部准备
func Ready(room Room) int {
	a := 0
	for _, item := range room.User {
		if item.Ready == 1 {
			a = a + 1
		}
	}
	if a == room.People {
		return 1
	} else {
		return 0
	}
}

// 用户准备
func Re(room Room, user string) Room {
	for _, item := range room.User {
		if item.OpenID == user {
			item.Ready = 1
			ServerSend(room, "用户"+item.OpenID+"已经准备")
			break
		}
	}
	return room
}

// 身份分配
func Iden(room Room, user string, iden string) Room {
	if room.GameMode == "新手场" {
		Ww := 2
		God := 1
		Solor := 1
		Peo := 2
		// 新手场抢身份
		for _, item := range room.User {
			if item.OpenID == user {
				switch iden {
				case "狼人":
					if Ww != 0 {
						Ww = Ww - 1
						item.Identity = "狼人"
					}
				case "预言家":
					if God != 0 {
						God = God - 1
						item.Identity = "预言家"
					}
				case "猎人":
					if Solor != 0 {
						Solor = Solor - 1
						item.Identity = "猎人"
					}
				default:

				}
			}
		}

	}
}

// 判断是否结束游戏
func Over(room Room) int {
	a := 0
	for _, item := range room.User {
		if item.Ready == 1 {
			a = a + 1
		}
	}
	if a == room.People {
		return 1
	} else {
		return 0
	}
}

// 游戏内容
func GameIng(room Room) {
	switch room.GameMode {
	case "新手场":
		a := 1
		for {
			ServerSend(room, "第"+strconv.Itoa(a)+"轮开始")
		}
	}
}
