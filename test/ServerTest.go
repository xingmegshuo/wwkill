/***************************
@File        : Server_test.go
@Time        : 2020/12/21 11:28:16
@AUTHOR      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : 服务端测试，模拟客户端
****************************/

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

//  服务端测试
func main() {
	conn, err := net.Dial("tcp", "localhost:4321")
	if err != nil {
		log.Println(err, "服务无法开启")
		return
	}
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("client start !!!")
	// fmt.Println("First, what is your name?")
	// clientName, _ := inputReader.ReadString('\n')
	// fmt.Printf("CLIENTNAME %s", clientName)
	// trimmedClient := strings.Trim(clientName, "\r\n") // Windows 平台下用 "\r\n"，Linux平台下使用 "\n"
	// 给服务器发送信息直到程序退出：
	for {

		fmt.Println("What to send to the server? Type Q to quit.")
		input, _ := inputReader.ReadString('\n')
		trimmedInput := strings.Trim(input, "\n")
		if trimmedInput == "Q" {
			return
		}
		_, err = conn.Write([]byte(trimmedInput))
		if err != nil {
			log.Println("服务端连接丢失", err.Error())
			continue
		}
		mes := make([]byte, 512)
		len, err := conn.Read(mes)
		if err != nil {
			log.Println("服务器丢失", err.Error())
			continue
		}
		if len > 0 {
			fmt.Println("服务端返回:", string(mes[:len]))
		}
	}
}
