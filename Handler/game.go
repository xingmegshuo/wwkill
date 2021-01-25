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
	"wwKill/Mydb"

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
	Score    int
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

// 结束游戏
func GameOver(room Room, over int) {
	for _, item := range room.User {
		thisUser := Mydb.User{
			OpenID: item.OpenID,
		}
		U, _ := ctrlUser.GetUser(thisUser)
		if over == 1 {
			if item.Identity == "狼人" {
				record := Mydb.Record{
					User:     U.Id,
					GameTime: time.Now(),
					Identity: item.Identity,
					GameMode: room.GameMode,
					RunAway:  0,
					Result:   "胜利",
				}
				ctrlRecord.Insert(record)
			} else {
				record := Mydb.Record{
					User:     U.Id,
					GameTime: time.Now(),
					Identity: item.Identity,
					GameMode: room.GameMode,
					RunAway:  0,
					Result:   "失败",
				}
				ctrlRecord.Insert(record)

			}
		} else {
			if item.Identity == "狼人" {
				record := Mydb.Record{
					User:     U.Id,
					GameTime: time.Now(),
					Identity: item.Identity,
					GameMode: room.GameMode,
					RunAway:  0,
					Result:   "失败",
				}
				ctrlRecord.Insert(record)
			} else {
				record := Mydb.Record{
					User:     U.Id,
					GameTime: time.Now(),
					Identity: item.Identity,
					GameMode: room.GameMode,
					RunAway:  0,
					Result:   "胜利",
				}
				ctrlRecord.Insert(record)

			}
		}
	}
}

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
			Wi:       0,
		}
	case "普通场":
		room = Room{
			GameMode: "普通场",
			People:   9,
			Public:   1,
			Ww:       3,
			God:      1,
			Ci:       3,
			Hu:       1,
			Wi:       1,
		}
	case "高手场":
		room = Room{
			GameMode: "高手场",
			People:   10,
			Public:   1,
			Ww:       3,
			Hu:       7,
			God:      0,
			Ci:       0,
			Wi:       0,
		}
	}
	return room
}

// 游戏逻辑
func GameSocket(room Room, ch chan string) {
	a := 1
	var ro Room
	ro = room
	wait := ""
	for str := range ch {
		switch str[:4] {
		case "died":
			// to do 用户死亡
			ro = Die(room, str[4:])
		case "save":
			// 救活用户
			ro = Save(room, str[4:])
		case "diew":
			wait = WaitSave(room, str[4:])
		}
	}

	for {
		over := Over(ro)
		if over == 1 {
			ServerSend(ro, "游戏结束,狼人胜利")
			GameOver(ro, over)
			break
		}
		if over == 2 {
			ServerSend(ro, "游戏结束,平民胜利")
			GameOver(ro, over)
			break
		}
		ServerSend(ro, "法官:第"+strconv.Itoa(a)+"天")
		time.Sleep(time.Second * 2)
		Black(room, strconv.Itoa(a), wait)
		ServerSend(ro, "第"+strconv.Itoa(a)+"天:天亮了请睁眼")
		Day(room)
	}
}

// 加入房间
func Join(room Room, player Player) Room {
	room.User = append(room.User, player)
	ServerSend(room, "用户"+player.OpenID+"进入房间")
	return room
}

// 房间连接
func RoomSocket(conn []byte) {
	var mes Mes
	err := json.Unmarshal(conn, &mes)
	if err != nil {
		log.Println("连接断开")
		// todo 退出
	} else {
		for _, room := range PlayRoom {
			ch := make(chan string)
			if room.Owner == mes.Room {
				if mes.Message == "准备" {
					room = Re(room, mes.User)
				}
				if mes.Message[:2] == "身份" {
					// to do 分配身份
					room = Iden(room, mes.User, mes.Message[2:])
					GameSocket(room, ch)
				}
				if len(room.User) == room.People {
					ServerSend(room, "法官:人员已到齐，请所有人准备!")
					continue
				}
				if mes.Message[:2] == "查看" {
					// ServerSend(room, mes.User+"查看"+mes.Message[2:])
					// 预言家查看身份
					LookIden(mes.User, room, mes.Message[:2])
				}
				if mes.Message[:2] == "毒药" {
					WiKill(mes.User, room, mes.Message[2:], ch)
				}
				if mes.Message[:2] == "解药" {
					// 女巫救人
					WiSave(mes.User, room, mes.Message[2:], ch)
				}
				if mes.Message[:2] == "暗杀" {
					// 狼人杀人
					WwKill(mes.User, room, mes.Message[2:], ch)
				}
				if mes.Message[:2] == "杀人" {
					// 猎人杀人
					HuKill(mes.User, room, mes.Message[2:], ch)
				}
				if mes.Message[:2] == "投票" {
					// 大家投票
					HuKill(mes.User, room, mes.Message[2:], ch)
				}
				if mes.Message[:2] == "离开" {
					// 退出房间
					room = Leave(mes.User, room)
				}
				ready := Ready(room)
				if ready == 1 {
					ServerSend(room, "法官:所有人已准备，游戏5秒后开始!")
					time.Sleep(time.Second * 5)
					ServerSend(room, "游戏开始！请选择身份")
				}
			}
		}
	}
}

// 退出房间
func Leave(user string, room Room) Room {
	a := 0
	for l, item := range room.User {
		if item.OpenID == user {
			a = l
			// openId 赋值给游戏中用户
			client_user[item.Ws] = client_palyer[item.Ws]
			// 把用户从在线用户移除
			delete(client_palyer, item.Ws)
		}
	}
	room.User = room.User[:a+copy(room.User[a:], room.User[a+1:])]
	return room
}

// 等待死亡
func WaitSave(room Room, user string) string {
	mes := ""
	for _, item := range room.User {
		if item.OpenID == user {
			mes = item.OpenID
			break
		}
	}
	return mes
}

// 用户死亡
func Die(room Room, user string) Room {
	for _, item := range room.User {
		if item.OpenID == user {
			item.Survive = 3
			break
		}
	}
	return room
}

// 救活用户
func Save(room Room, user string) Room {
	for _, item := range room.User {
		if item.OpenID == user {
			item.Survive = 1
			break
		}
	}
	return room
}

// 猎人杀人
func HuKill(user string, room Room, look string, ch chan string) {
	for _, item := range room.User {
		if item.OpenID == look {
			ch <- "died" + item.OpenID
		}
		if item.OpenID == user {
			Send(item.Ws, "您开枪带走了"+look)
		}
	}
}

// 狼人杀人
func WwKill(user string, room Room, look string, ch chan string) {
	score := 0
	kill := ""
	for _, item := range room.User {
		if item.OpenID == look {
			item.Score = item.Score + 1
		}
		if item.OpenID == user {
			Send(item.Ws, "您投票给"+look)
		}
		if item.Score > score {
			score = item.Score
			kill = item.OpenID
		}
	}
	if room.Wi > 0 {
		ch <- "diew" + kill
	} else {
		ch <- "died" + kill
	}
}

// 女巫救人
func WiSave(user string, room Room, look string, ch chan string) {
	for _, item := range room.User {
		if item.OpenID == look {
			ch <- "save" + item.OpenID
		}
		if item.OpenID == user {
			Send(item.Ws, "您用解药救了"+look)
		}
	}
}

// 预言家查看身份
func LookIden(user string, room Room, look string) {
	iden := ""
	for _, item := range room.User {
		if item.OpenID == look {
			iden = item.Identity
			continue
		}
		if item.OpenID == user && iden != "" {
			Send(item.Ws, "您查看了"+look+"它的身份是"+iden)
		}
	}
}

// 女巫毒人
func WiKill(user string, room Room, look string, ch chan string) {
	for _, item := range room.User {
		if item.OpenID == look {
			ch <- "died" + item.OpenID
		}
		if item.OpenID == user {
			Send(item.Ws, "您用毒药毒死了"+look)
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
		OpenID:  client_palyer[ws],
		Ws:      ws,
		Survive: 1,
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

// 发送给狼人
func ServerWw(room Room, str string) {
	for _, U := range room.User {
		if U.Identity == "狼人" && U.Survive == 1 {
			Send(U.Ws, str)
		}
	}
}

// 发送给预言家
func ServerGod(room Room, str string) {
	for _, U := range room.User {
		if U.Identity == "预言家" && U.Survive == 1 {
			Send(U.Ws, str)
		}
	}
}

// 发送给女巫
func ServerWi(room Room, str string) {
	for _, U := range room.User {
		if U.Identity == "女巫" && U.Survive == 1 {
			Send(U.Ws, str)
		}
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
	for _, item := range room.User {
		if item.OpenID == user {
			switch iden {
			case "狼人":
				if room.Ww != 0 {
					room.Ww = room.Ww - 1
					item.Identity = "狼人"
				}
			case "预言家":
				if room.God != 0 {
					room.God = room.God - 1
					item.Identity = "预言家"
				}
			case "猎人":
				if room.Hu != 0 {
					room.Hu = room.Hu - 1
					item.Identity = "猎人"
				}
			case "女巫":
				if room.Wi != 0 {
					room.Wi = room.Wi - 1
					item.Identity = "女巫"
				}
			default:
				idenstr := ""
				if room.Ww > 0 {
					room.Ww = room.Ww - 1
					idenstr = "狼人"
				}
				if room.God > 0 {
					room.Ww = room.Ww - 1
					idenstr = "预言家"
				}
				if room.Hu > 0 {
					room.Hu = room.Hu - 1
					idenstr = "猎人"
				}
				if room.Wi > 0 {
					room.Wi = room.Wi - 1
					idenstr = "女巫"
				}
				if room.Ci > 0 {
					room.Ci = room.Ci - 1
					idenstr = "平民"
				}
				item.Identity = idenstr
			}
			Send(item.Ws, "您的身份是:"+item.Identity)
		}
	}
	return room
}

// 判断是否结束游戏
func Over(room Room) int {
	a := 0
	b := 0
	c := 0
	for _, item := range room.User {
		if item.Survive == 0 {
			if item.Identity == "狼人" {
				a = a + 1
			}
			if item.Identity == "平民" {
				b = b + 1
			}
			if item.Identity == "女巫" || item.Identity == "猎人" || item.Identity == "预言家" {
				c = c + 1
			}

		}
	}
	switch room.GameMode {
	case "新手场":
		if a == room.Ww {
			return 1
		}
		if b == room.Ci || c == room.God+room.Hu+room.Wi {
			return 2
		}
	case "普通场":
		if a == room.Ww {
			return 1
		}
		if b == room.Ci || c == room.God+room.Hu+room.Wi {
			return 2
		}
	case "高手场":
		if a == room.Ww {
			return 1
		}
		if c == room.God+room.Hu+room.Wi {
			return 2
		}
	}

	return 0
}

// 天黑阶段
func Black(room Room, day string, wait string) {
	ServerSend(room, "法官:天黑了")
	time.Sleep(time.Second * 3)
	ServerWw(room, "请狼人开始行动")
	if room.God > 0 {
		ServerGod(room, "请预言家查验身份")
	}
	time.Sleep(time.Second * 20)
	if room.Wi > 0 {
		ServerWi(room, "请女巫选择是否用药")
		if wait != "" {
			ServerWi(room, "昨晚上被杀的是"+wait+"是否救一下")
		}
	}
	time.Sleep(time.Second * 20)
}

// 白天阶段
func Day(room Room) {
	ServerSend(room, "法官:天亮了")
	for _, item := range room.User {
		if item.Survive == 3 {
			item.Survive = 0
			if item.Identity == "猎人" {
				ServerSend(room, "法官:用户"+item.OpenID+"死亡,他的身份是猎人请他发动技能")
				time.Sleep(time.Second * 10)
			} else {
				ServerSend(room, "法官:用户"+item.OpenID+"死亡,请他发言")
				time.Sleep(time.Second * 30)
			}
		}
		if item.Survive == 1 {
			item.Survive = 0
			ServerSend(room, "法官:请用户"+item.OpenID+"发言")
		}
	}
	ServerSend(room, "法官:请用户投票")
	time.Sleep(time.Second * 15)
	for _, item := range room.User {
		if item.Survive == 3 {
			item.Survive = 0
			if item.Identity == "猎人" {
				ServerSend(room, "法官:用户"+item.OpenID+"死亡,他的身份是猎人请他发动技能")
				time.Sleep(time.Second * 10)
			}
			ServerSend(room, "法官:用户"+item.OpenID+"死亡,请他发言")
			time.Sleep(time.Second * 30)
		}
		if item.Survive == 2 {
			ServerSend(room, "法官:请用户"+item.OpenID+"发言")
			item.Survive = 0
		}
	}
}
