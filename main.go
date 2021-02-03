/***************************
@File        : main.go
@Time        : 2020/12/23 13:05:53
@AUTHOR      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : 运行主程序
****************************/

package main

import (
	"log"
	"net/http"
	"wwKill/Handler"

	"golang.org/x/net/websocket"
	)

var client_map = make(map[*websocket.Conn]string)

// 开启服务端
func main() {
	log.Println("服务开启")
	http.Handle("/", websocket.Handler(Echo))

	if err := http.ListenAndServe(":4321", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// websocket
func Echo(ws *websocket.Conn) {
	// websocket.Conn  用来作为客户端和服务器端交互的通道
	// 只是用来记录接收请求的次数
	// 一直接收连接

	ParseConn(ws)
}

// 解析接收的数据
func ParseConn(ws *websocket.Conn) {
	var err error
	client_map[ws] = ws.RemoteAddr().String()
	log.Println(client_map)
	for {
		var reply string
		if err = websocket.Message.Receive(ws, &reply); err != nil {
			log.Println("客户端丢失")
			CloseConn(ws)
			break
		}
		log.Printf("用户发送了: %v\n", string(reply))

		// 发送给解析函数
		go Handler.ParseData(string(reply), ws)
		// log.Println("等待三秒返回测试")
		// time.Sleep(time.Duration(20) * time.Second)

	}
}

// 断开连接
func CloseConn(ws *websocket.Conn) {
	delete(client_map, ws)
	Handler.CloseUser(ws)
	ws.Close()
	// log.Println(len(client_map))
}
