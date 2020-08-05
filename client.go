package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)


type MyInfo struct {
	Room Room
	Name string
}

type Room struct {
	Name string
}

type Msg struct {
	Type string
	Text string
	Info MyInfo
}


func InClientName() MyInfo {
	var s string
	fmt.Print("사용자 이름 입력 : ")
	fmt.Scanln(&s)
	var r = Room {}
	var mi = MyInfo{r, s}
	return mi
}

func init() {
	fmt.Println("init")
}

func RecvServer(conn net.Conn) {
	defer func(){
		if r:= recover(); r != nil {
			log.Fatalln(r)
		}
	}()

	var recvBuf []byte
	recvBuf = make([]byte, 4096)

	for {
		n, err := conn.Read(recvBuf)
		if err != nil {
			log.Println("EOF", err)
			panic("close")
		}
		data := recvBuf[:n]
		var m Msg
		json.Unmarshal([]byte(data), &m)

		log.Println("Server send: ", m)
	}
}


func SendServerMag(conn net.Conn, mi MyInfo){

	defer func(){
		if r:= recover(); r != nil {
			log.Fatalln(r)
		}
	}()

	for {
		var s string
		var t string
		// 방 만들기 + 방 입장
		if mi.Room.Name == "" {
			t = "room"
			fmt.Scanln(&s)
			var room = Room{string(s)}
			mi.Room = room
		} else {
			t = "text"
			fmt.Scanln(&s)

			if s == "exit" {
				log.Fatalln("exit")
			}
		}

		var m = Msg {t, string(s), mi}
		b, _ := json.Marshal(m)
		conn.Write(b)
	}
}


func main() {
	conn, err := net.Dial("tcp", ":8000")
	if nil != err {
		log.Println("net.Dial", err)
	}

	defer func(){
		r := recover()
		log.Println("exit", r)
	}()


	var mi = InClientName()

	go RecvServer(conn)
	go SendServerMag(conn, mi)

	select{}

}