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
	"net"
	"time"
	"wwKill/Handler"
)

// 开启服务端
func main() {
	listener, err := net.Listen("tcp", "localhost:4321")
	if err != nil {
		log.Println(err.Error(), "服务无法开启")
		return
	}
	log.Println("服务开启!!!")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting", err.Error())
		}
		log.Println("收到用户连接")
		go ParseConn(conn)
	}
}

// 解析接收的数据
func ParseConn(conn net.Conn) {
	for {
		log.Println("监听收到的用户数据")
		buf := make([]byte, 512)
		len, err := conn.Read(buf)
		if err != nil {
			log.Println("Error reading", err.Error())
			break
		}
		log.Printf("用户发送了: %v\n", string(buf[:len]))
		// 发送给解析函数
		mes := Handler.ParseData(string(buf[:len]))
		log.Println("等待三秒返回测试")
		time.Sleep(time.Duration(3) * time.Second)
		_, err = conn.Write([]byte(mes))
		if err != nil {
			log.Println("客户端丢失", err.Error())
			continue
		}
	}
}
