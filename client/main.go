package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func MessageSend(conn net.Conn) {
	var input string
	//不断读取终端输入的字符串
	for {
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		input = string(data)
		if strings.ToUpper(input) == "EXIT" {
			conn.Close()
			break
		}
		_, err := conn.Write([]byte(input))
		if err != nil {
			conn.Close()
			fmt.Println("客户端连接错误:" + err.Error())
			break
		}
	}
}
func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	CheckError(err)
	defer conn.Close()
	//conn.Write([]byte("hello golang"))
	//处理用户输入
	go MessageSend(conn)
	buf := make([]byte, 1024)
	for {
		length, err := conn.Read(buf)
		CheckError(err)
		fmt.Println("receive server message content:" + string(buf[:length]))
	}
	fmt.Println("Client program end!")
}
