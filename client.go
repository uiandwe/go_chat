package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"time"
)

type MsgBody struct {
	Content string
}

type Msg struct {
	Header MsgHeader
	Body   interface{}
}
type MsgHeader struct {
	MsgType string
	Date    string
}

func init() {
	fmt.Println("init")
	gob.Register(MsgBody{})

}

func RecvServer(conn net.Conn) {
	defer func(){
		if r:= recover(); r != nil {
			log.Fatalln(r)
		}
	}()

	data := make([]byte, 4096)

	for {
		n, err := conn.Read(data)
		if err != nil {
			log.Println("EOF", err)
			panic("close")
		}

		log.Println("Server send : " + string(data[:n]))
	}
}


func SendServerMag(conn net.Conn){

	defer func(){
		if r:= recover(); r != nil {
			log.Fatalln(r)
		}
	}()

	var (
		codeBuffer bytes.Buffer
		enc         *gob.Encoder = gob.NewEncoder(&codeBuffer)
	)

	for {
		var s string
		fmt.Scanln(&s)
		if s == "exit" {
			log.Fatalln("exit")
		}

		enc.Encode(Msg{
			Header: MsgHeader{
				MsgType: "text",
				Date:    time.Now().UTC().Format(time.RFC3339),
			},
			Body: MsgBody{
				Content: string(s),
			},
		})

		conn.Write(codeBuffer.Bytes())
		codeBuffer.Reset()
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