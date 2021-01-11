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

// 开启服务端
func main() {
	log.Println("服务开启")
	http.Handle("/", websocket.Handler(Echo))

	if err := http.ListenAndServe(":4321", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
	// 访问服务器的地址,ip没有限制,端口是8888,
	//  ws://127.0.0.1:8888

	// listener, err := net.Listen("tcp", "localhost:4321")
	// if err != nil {
	// 	log.Println(err.Error(), "服务无法开启")
	// 	return
	// }
	// log.Println("服务开启!!!")
	// for {
	// 	conn, err := listener.Accept()
	// 	if err != nil {
	// 		log.Println("Error accepting", err.Error())
	// 	}
	// 	log.Println("收到用户连接")
	// 	go ParseConn(conn)
	// }
}

// websocket
func Echo(ws *websocket.Conn) {
	// websocket.Conn  用来作为客户端和服务器端交互的通道
	// 只是用来记录接收请求的次数
	ParseConn(ws)
	// for {
	// 	log.Println("连接")
	// 	var reply string
	// 	// 建立连接后 接收来自客户端的信息reply
	// 	if err = websocket.Message.Receive(ws, &reply); err != nil {
	// 		fmt.Println("客户端丢失")
	// 	}
	// 	log.Println("接收的信息: " + reply)
	// 	// 把收到的信息进行处理,也可以做信息过滤,也可以返回固定的信息
	// 	msg := "Received:  " + reply
	// 	log.Println("发给客户端的信息: " + msg)
	// 	// 把信息返回发送给客户端
	// 	if err = websocket.Message.Send(ws, msg); err != nil {
	// 		log.Println("客户端丢失")
	// 	}
	// }
}

// 解析接收的数据
func ParseConn(ws *websocket.Conn) {
	var err error
	for {
		log.Println("用户连接")
		var reply string
		if err = websocket.Message.Receive(ws, &reply); err != nil {
			log.Println("客户端丢失")
			break
		}
		log.Printf("用户发送了: %v\n", string(reply))
		// 发送给解析函数
		mes := Handler.ParseData(string(reply))
		// log.Println("等待三秒返回测试")
		// time.Sleep(time.Duration(20) * time.Second)
		if err = websocket.Message.Send(ws, mes); err != nil {
			log.Println("客户端丢失", err.Error())
			continue
		}
	}
}
