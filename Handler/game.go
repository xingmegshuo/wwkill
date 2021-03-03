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
	Status   string
}

type Room struct {
	GameMode string
	People   int
	Public   int
	Stop     int
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

// 清除无效用户
func ClearUser(){
	for _, ro := range PlayRoom{
		for _, u := range ro.User{
			if _,ok:= client_palyer[u.Ws];!ok{
				log.Println("清除用户",u.OpenID)
				Remove(u.OpenID)
			}
		}
	}
}


// 开始游戏
func GameStart(mes []byte, ws *websocket.Conn) string {
	var room Room
	ClearUser()
	err := json.Unmarshal(mes, &Game)
	if err != nil {
		log.Println("数据问题:", err.Error())
	}
	if Game.New == 0 {
		room = SearchRoom(Game.GameMode)
	} else {
		room = NewRoom(Game.GameMode)
	}
	log.Println(len(client_palyer),"几个玩家",len(room.User))
	message := Init(ws, room)
	log.Println(len(client_palyer),"几个玩家","房间号",message)
	return ToMes("roomID", message)
}

// 查找房间
func SearchRoom(GameType string) Room {
	var room Room
	new := "true"
	for _, v := range PlayRoom {
		if len(v.User) < v.People && v.Public == 1 && v.GameMode == GameType && v.Stop == 0 {
			room = v
			new = "false"
			break
		}
	}
	if new == "true" || len(PlayRoom) == 0 {
		room = NewRoom(GameType)
	}
	return room

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
			Stop:     0,
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
			Stop:     0,
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
			Stop:     0,
			Ww:       3,
			Hu:       7,
			God:      0,
			Ci:       0,
			Wi:       0,
		}
	}
	return room
}

// 加入房间
func Join(room Room, player Player) Room {
	add := 0
	for l, ro := range room.User {
		if ro.OpenID == player.OpenID {
			add = 1
			room.User[l] = player
		}
	}
	if add == 0 {
		room.User = append(room.User, player)
	}
	ServerSend(room, player.OpenID+":用户"+player.OpenID+"进入房间")
	Update(room)
	return room
}

// 获取房间号
func GetRoomID(room Room) string {
	for l, ro := range PlayRoom {
		// log.Println(ro.Owner)
		if ro.Owner == room.Owner {
			// log.Println("hhhhhhhhhhhhhhh-房间号")
			return l
		}
	}
	return "null"
}

// 初始化改变用户状态,由在线换为游戏中
func Init(ws *websocket.Conn, room Room) string {
	// openId 赋值给游戏中用户
	client_palyer[ws] = client_user[ws]
	// 把用户从在线用户移除
	delete(client_user, ws)
	player := Player{
		OpenID:  client_palyer[ws],
		Ws:      ws,
		Survive: 1,
		Status:  "true",
	}
	if len(room.User) == 0 {
		room.Owner = player.OpenID
		PlayRoom[room.Owner] = room
		room = Join(room, player)
	} else {
		room = Join(room, player)
	}
	str := ""
	for l, item := range room.User {
		if l == len(room.User)-1 {
			str = str + "" + item.OpenID + ":" + strconv.Itoa(item.Ready) + ""
		} else {
			str = str + "" + item.OpenID + ":" + strconv.Itoa(item.Ready) + ","
		}
	}
	str = str + ":房间总人数"
	ServerSend(room, str)
	log.Println(len(room.User),"用户",len(PlayRoom),"房间个数","房间问题")
	return GetRoomID(room)
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
					User:     int(U.Id),
					GameTime: time.Now().Format("2006-01-02 15:04:05"),
					Identity: item.Identity,
					GameMode: room.GameMode,
					RunAway:  0,
					Result:   "胜利",
				}
				ctrlRecord.Insert(record)
			} else {
				record := Mydb.Record{
					User:     int(U.Id),
					GameTime: time.Now().Format("2006-01-02 15:04:05"),
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
					User:     int(U.Id),
					GameTime: time.Now().Format("2006-01-02 15:04:05"),
					Identity: item.Identity,
					GameMode: room.GameMode,
					RunAway:  0,
					Result:   "失败",
				}
				ctrlRecord.Insert(record)
			} else {
				record := Mydb.Record{
					User:     int(U.Id),
					GameTime: time.Now().Format("2006-01-02 15:04:05"),
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

// 房间连接
func RoomSocket(conn []byte) {
	var value Mes
	err := json.Unmarshal(conn, &value)
	if err != nil {
		log.Println("连接断开", err)
		// todo 退出
	} else {
		for _, room := range PlayRoom {
			IsOutline(room)
			if room.Owner == value.Room {
				sock := 1
				if value.Message[:6] == "准备" {
					Re(room, value.User)
					ready := Ready(room)
					if ready == 1 {
						ServerSend(room, "法官:所有人已准备，游戏5秒后开始!")
						time.Sleep(time.Second * 5)
						ServerSend(room, "游戏开始！请选择身份")
					}
				}
				if value.Message[:6] == "身份" {
					// to do 分配身份
					Iden(room, value.User, value.Message[6:])
				}
				if value.Message[:6] == "离开" {
					// 退出房间
					Leave(value.User, room)
					sock = 0
				}
				if value.Message[:6] == "查验" {
					// 预言家查看身份
					LookIden(value.User, room, value.Message[6:])
					sock = 0
				}
				if value.Message[:6] == "毒药" {
					WiKill(value.User, room, value.Message[6:])
					sock = 0
				}
				if value.Message[:6] == "解药" {
					// 女巫救人
					sock = 0
				}
				if value.Message[:6] == "暗杀" {
					// 狼人杀人
					WwKill(value.User, room, value.Message[6:])
					sock = 0
				}
				if value.Message[:6] == "杀人" {
					// 猎人杀人
					HuKill(value.User, room, value.Message[6:])
					sock = 0
				}
				if value.Message[:6] == "聊天" {
					// 猎人杀人
					GamingChat(value.User, room, value.Message[6:])
					sock = 0
				}
				if value.Message[:6] == "投票" {
					// 大家投票
					WwKill(value.User, room, value.Message[6:])
					sock = 0
				}
				if sock == 1 {
					Gaming(room)
				}
			}
		}
	}
}
// 聊天
func GamingChat(user string, room Room, str string){
	ServerSend(room, "聊天内容," + user+":"+str)
}

// 判断是否离线
func IsOutline(room Room) {
	for _, u := range room.User {
		if _, ok := client_palyer[u.Ws]; !ok {
			Remove(u.OpenID)
		}
	}
}

// 游戏中
func Gaming(room Room) {
	start := Start(room)
	if start == 0 {
		for _, a := range []int{1, 2, 3, 4, 5, 6, 7} {
			// Update(room)
			if a == 1 {
				ServerSend(room, "法官:start Game!!!!")
			}
			ServerSend(room, "法官:第"+strconv.Itoa(a)+"天")
			time.Sleep(time.Second * 2)
			Black(room, strconv.Itoa(a))
			ServerSend(room, "第"+strconv.Itoa(a)+"天:天亮了请睁眼")
			b := Day(room, a)
			if b == 1 {
				break
			}
		}
	}
}

// 判断是否结束
func IsOver(room Room) int {
	over := Over(room)
	if over == 1 {
		ServerSend(room, "游戏结束,平民胜利")
		GameOver(room, over)
		str := ""
		for l, item := range room.User {
			if l == len(room.User)-1 {
				str = str + "" + item.OpenID + ":" + item.Identity + ""
			} else {
				str = str + "" + item.OpenID + ":" + item.Identity + ","
			}
		}
		str = str + ":身份"
		ServerSend(room, str)
		return 1
	}
	if over == 2 {
		ServerSend(room, "游戏结束,狼人胜利")
		GameOver(room, over)
		str := ""
		for l, item := range room.User {
			if l == len(room.User)-1 {
				str = str + "" + item.OpenID + ":" + item.Identity + ""
			} else {
				str = str + "" + item.OpenID + ":" + item.Identity + ","
			}
		}
		str = str + ":身份"
		ServerSend(room, str)
		return 1
	}
	return 0
}

// 退出房间
func Leave(user string, room Room) {
	a := 0
	ServerSend(room, user+":退出房间")
	// 删除用户
	for i, u := range room.User {
		if u.OpenID == user {
			a = i
			client_user[u.Ws] = client_palyer[u.Ws]
			delete(client_palyer, u.Ws)
		}
	}
	room.User = append(room.User[:a], room.User[a+1:]...)
	Update(room)
}

// 等待死亡
func WaitSave(room Room, user string) {
	for l, item := range room.User {
		if item.OpenID == user {
			item.Survive = 2
			room.User[l] = item
			break
		}
	}
	Update(room)
}

// 更新
func Update(room Room) {
	for i, ro := range PlayRoom {
		if ro.Owner == room.Owner {
			if len(room.User) ==0{
				delete(PlayRoom, i)
			}else{
				PlayRoom[i] = room
			}
		}
	}
}

// 用户死亡
func Die(room Room, user string) {
	for l, item := range room.User {
		if item.OpenID == user {
			item.Survive = 3
			room.User[l] = item
			break
		}
	}
	Update(room)
}

// 枪杀
func HuDie(room Room, user string) {
	for l, item := range room.User {
		if item.OpenID == user {
			item.Survive = 4
			room.User[l] = item
			break
		}
	}
	Update(room)
}

// 救活用户
func Save(room Room, user string) {
	for l, item := range room.User {
		if item.OpenID == user {
			item.Survive = 1
			room.User[l] = item
			break
		}
	}
	Update(room)
}

// 猎人杀人
func HuKill(user string, room Room, look string) {
	for _, item := range room.User {
		if item.OpenID == user {
			HuDie(room, look)
			SendMS(item.Ws, "您开枪带走了"+look)
		}
	}
}

// 狼人杀人
func WwKill(user string, room Room, look string) {
	for l, item := range room.User {
		if item.OpenID == look {
			item.Score = item.Score + 1
		}
		if item.OpenID == user {
			SendMS(item.Ws, "您投票给"+look)
		}
		room.User[l] = item
	}
	Update(room)
}

// 女巫救人
func WiSave(user string, room Room, look string) {
	for _, item := range room.User {
		if item.OpenID == user {
			Save(room, look)
			SendMS(item.Ws, "您用解药救了"+look)
		}
	}
}

// 预言家查看身份
func LookIden(user string, room Room, look string) {
	iden := ""
	l := 0
	for i, item := range room.User {
		if item.OpenID == look {
			iden = item.Identity
		}
		if item.OpenID == user {
			l = i
		}
	}
	SendMS(room.User[l].Ws, "您查看了"+look+"它的身份是"+iden)
}

// 女巫毒人
func WiKill(user string, room Room, look string) {
	for _, item := range room.User {
		if item.OpenID == user {
			Die(room, look)
			SendMS(item.Ws, "您用毒药毒死了"+look)
		}
	}
}

// 房间内设置掉线
func Remove(str string) {
	for _, item := range PlayRoom {
		a := -1
		for n, user := range item.User {
			if user.OpenID == str && item.Stop == 1 {
				delete(client_palyer, user.Ws)
				user.Status = "false"
				item.User[n] = user
			}
			if user.OpenID == str && item.Stop == 0{
				a = n
			}
		}
		if a != -1{
			item.User = append(item.User[:a], item.User[a+1:]...)
		}
		log.Println(len(item.User),"房间人数---------------------")
		Update(item)
	}
}

// 房间内设置掉线
func Close(ws *websocket.Conn) {
	for _,ro := range PlayRoom{
		for _, u := range ro.User{
			if u.Ws == ws{
				Remove(u.OpenID)
			}
		}
	}
	// Remove(client_palyer[ws])
}

// 发送服务器公告信息
func ServerSend(room Room, str string) {
	for _, U := range room.User {
		if U.Status == "true" {
			SendMS(U.Ws, str)
		}
	}
}

// 发送给狼人
func ServerWw(room Room, str string) {
	for _, U := range room.User {
		if U.Identity == "狼人" && U.Survive == 1 && U.Status == "true" {
			SendMS(U.Ws, str)
		}
	}
}

// 发送给预言家
func ServerGod(room Room, str string) {
	for _, U := range room.User {
		if U.Identity == "预言家" && U.Survive == 1 && U.Status == "true" {
			SendMS(U.Ws, str)
		}
	}
}

// 发送给女巫
func ServerWi(room Room, str string) {
	for _, U := range room.User {
		if U.Identity == "女巫" && U.Survive == 1 && U.Status == "true" {
			SendMS(U.Ws, str)
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
		for l, item := range PlayRoom {
			if item.Owner == room.Owner {
				item.Stop = 1
				PlayRoom[l] = item
			}
		}
		return 1
	} else {
		return 0
	}
}

// 用户准备
func Re(room Room, user string) {
	for l, item := range room.User {
		if item.OpenID == user {
			if item.Ready == 0 {
				item.Ready = 1
				ServerSend(room, item.OpenID+":用户"+item.OpenID+"已经准备")
			} else {
				item.Ready = 0
				ServerSend(room, item.OpenID+":用户"+item.OpenID+"取消准备")
			}
			room.User[l] = item
		}
	}
	Update(room)
}

// 游戏开始
func Start(room Room) int {
	start := 0
	for _, item := range room.User {
		if item.Identity == "" {
			start = 1
		}
	}
	return start
}

// 身份分配
func Iden(room Room, user string, iden string) {
	for l, item := range room.User {
		if item.OpenID == user {
			switch iden {
			case "狼人":
				if room.Ww != 0 {
					room.Ww = room.Ww - 1
					item.Identity = "狼人"
				} else {
					item.Identity = ""
				}
			case "预言家":
				if room.God != 0 {
					room.God = room.God - 1
					item.Identity = "预言家"
				} else {
					item.Identity = ""
				}
			case "猎人":
				if room.Hu != 0 {
					room.Hu = room.Hu - 1
					item.Identity = "猎人"
				} else {
					item.Identity = ""
				}
			case "女巫":
				if room.Wi != 0 {
					room.Wi = room.Wi - 1
					item.Identity = "女巫"
				} else {
					item.Identity = ""
				}
			default:
				idenstr := ""
				if room.Ww > 0 {
					room.Ww = room.Ww - 1
					idenstr = "狼人"
					item.Identity = idenstr

					break
				}
				if room.God > 0 {
					room.God = room.God - 1
					idenstr = "预言家"
					item.Identity = idenstr
					break
				}
				if room.Hu > 0 {
					room.Hu = room.Hu - 1
					idenstr = "猎人"
					item.Identity = idenstr
					break
				}
				if room.Wi > 0 {
					room.Wi = room.Wi - 1
					idenstr = "女巫"
					item.Identity = idenstr
					break
				}
				if room.Ci > 0 {
					room.Ci = room.Ci - 1
					idenstr = "平民"
					item.Identity = idenstr
					break
				}
			}
			// log.Println(room.Owner, "-------------", item.Identity, user)
			SendMS(item.Ws, "您的身份是:"+item.Identity)
		}
		room.User[l] = item
	}
	Update(room)

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
		if a == 2 {
			return 1
		}
		if b == 2 || c == 2 {
			return 2
		}
	case "普通场":
		if a == 3 {
			return 1
		}
		if b == 3 || c == 3 {
			return 2
		}
	case "高手场":
		if a == 3 {
			return 1
		}
		if c == 7 {
			return 2
		}
	}
	return 0
}

// 投票结果
func Result(room Room) {
	score := 0
	kill := ""
	wait := 0
	for _, item := range room.User {
		if item.Score > score && item.Survive != 0 {
			score = item.Score
			kill = item.OpenID
		}
		if item.Identity == "女巫" && item.Survive != 0 {
			wait = 1
		}
	}
	for _, item := range room.User {
		item.Score = 0
	}
	Update(room)
	if score != 0 {
		if wait == 1 {
			WaitSave(room, kill)
		} else {
			Die(room, kill)
		}
	}
}

// 天黑阶段
func Black(room Room, day string) {
	ServerSend(room, "法官:天黑了")
	time.Sleep(time.Second * 3)
	for _, item := range room.User {
		if item.Survive != 0 && item.Identity == "预言家" {
			ServerGod(room, "请预言家查验身份")
		}
	}
	ServerWw(room, "请狼人开始行动")
	time.Sleep(time.Second * 20)
	// 统计狼人投票结果
	Result(room)
	wait := ""
	for _, item := range room.User {
		if item.Survive == 2 {
			wait = item.OpenID
		}
		if item.Survive != 0 && item.Identity == "女巫" {
			ServerWi(room, "请女巫选择是否用药")
			time.Sleep(time.Second * 20)
			if wait != "" {
				ServerWi(room, "昨晚上被杀的是"+wait+"是否救一下")
			}
		}
	}
}

// 白天阶段
func Day(room Room, a int) int {
	ServerSend(room, "法官:天亮了")
	time.Sleep(time.Second * 2)
	for l, item := range room.User {
		if item.Survive == 3 {
			ServerSend(room, "法官:死亡用户,"+item.OpenID)
			time.Sleep(time.Second * 5)
			item.Survive = 0
			room.User[l] = item
			if item.Identity == "猎人" {
				ServerSend(room, "法官:用户"+item.OpenID+"死亡,他的身份是猎人请他发动技能")
				time.Sleep(time.Second * 10)
				if a == 1 {
					ServerSend(room, "法官:用户"+item.OpenID+"死亡,请他发言")
					time.Sleep(time.Second * 30)
				}
			} else {
				if a == 1 {
					ServerSend(room, "法官:用户"+item.OpenID+"死亡,请他发言")
					time.Sleep(time.Second * 5)
				}
			}
		}
	}
	Update(room)
	for l, item := range room.User {
		if item.Survive == 4 {
			item.Survive = 0
			ServerSend(room, "法官:死亡用户,"+item.OpenID)
			room.User[l] = item
		}
	}
	b := IsOver(room)
	if b == 1 {
		return 1
	}
	for _, item := range room.User {
		if item.Survive == 1 {
			item.Survive = 0
			ServerSend(room, "法官:请用户"+item.OpenID+"发言")
			time.Sleep(time.Second * 5)
		}
	}
	room = Clear(room)
	ServerSend(room, "法官:请用户投票")
	time.Sleep(time.Second * 15)
	Result(room)
	time.Sleep(time.Second * 2)
	for l, item := range room.User {
		if item.Survive == 3 {
			ServerSend(room, "法官:死亡用户,"+item.OpenID)
			time.Sleep(time.Second * 5)
			item.Survive = 0
			room.User[l] = item
			if item.Identity == "猎人" {
				ServerSend(room, "法官:用户"+item.OpenID+"死亡,他的身份是猎人请他发动技能")
				time.Sleep(time.Second * 10)
			}
			ServerSend(room, "法官:用户"+item.OpenID+"死亡,请他发言")
			time.Sleep(time.Second * 5)
		}
	}
	Update(room)
	for l, item := range room.User {
		if item.Survive == 4 {
			item.Survive = 0
			ServerSend(room, "法官:死亡用户,"+item.OpenID)
			room.User[l] = item
		}
	}
	Update(room)
	b = IsOver(room)
	if b == 1 {
		return 1
	}
	return 0
}

// 发送消息
func SendMS(ws *websocket.Conn, mes string) {
	if err := websocket.Message.Send(ws, mes); err != nil {
		log.Println("客户端丢失", err.Error())
		Close(ws)
		// log.Println("移除用户")
	}
}

// 清空投票状态
func Clear(room Room) Room {
	for l, item := range room.User {
		item.Score = 0
		room.User[l] = item
	}
	return room
}
