package main

import (
	"fmt"
	"net"
	"strings"
)

var (
	onlineConns  = make(map[string]net.Conn)
	messageQueue = make(chan string, 10000)
	quitChan     = make(chan bool)
)

//chatroom server
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func ProcessInfo(conn net.Conn) {
	buf := make([]byte, 1024)
	defer func(conn net.Conn) {
		addr := fmt.Sprintf("%s", conn.RemoteAddr())
		delete(onlineConns, addr)
		conn.Close()
		for conn := range onlineConns {
			fmt.Println(conn)
		}
	}(conn)
	for {
		numOfBytes, err := conn.Read(buf)
		if err != nil {
			break
		}
		if numOfBytes != 0 {
			message := string(buf[:numOfBytes])
			messageQueue <- message
		}
	}
}

func doProcessMessage(message string) {
	contents := strings.Split(message, "#")
	if len(contents) > 1 {
		addr := contents[0]
		sendMessage := strings.Join(contents[1:], "#")
		addr = strings.Trim(addr, " ")
		if conn, ok := onlineConns[addr]; ok {
			_, err := conn.Write([]byte(sendMessage))
			if err != nil {
				fmt.Println("online conns send failure!")
			}
		}
	}
}

func ComsumeMessage() {
	for {
		select {
		case message := <-messageQueue:
			//对消息进行解析
			doProcessMessage(message)
		case <-quitChan:
			break
		}
	}
}

func main() {
	onlineConns = make(map[string]net.Conn)
	listen, err := net.Listen("tcp", "127.0.0.1:8080")
	CheckError(err)
	defer listen.Close()
	fmt.Println("server is waitting......")
	go ComsumeMessage()
	for {
		conn, err := listen.Accept()
		CheckError(err)
		//将conn存储到onlineConns映射表中
		addr := fmt.Sprintf("%s", conn.RemoteAddr())
		onlineConns[addr] = conn
		for conn := range onlineConns {
			fmt.Println(conn)
		}
		go ProcessInfo(conn)
	}
}
