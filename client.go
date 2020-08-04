package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type Msg struct {
	Type string
	Text string
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


func SendServerMag(conn net.Conn){

	defer func(){
		if r:= recover(); r != nil {
			log.Fatalln(r)
		}
	}()

	for {
		var s string
		fmt.Scanln(&s)
		if s == "exit" {
			log.Fatalln("exit")
		}
		conn.Write([]byte(s))
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


	go RecvServer(conn)
	go SendServerMag(conn)

	select{}

}