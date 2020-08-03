package main

import (
	"fmt"
	"log"
	"net"
)


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

		log.Println("Server send: ", string(data))
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