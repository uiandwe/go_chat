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

	var codeBuffer bytes.Buffer
	var dec        *gob.Decoder = gob.NewDecoder(&codeBuffer)
	recvBuf := make([]byte, 4096)

	for {

		n, err := conn.Read(recvBuf)
		if err != nil {
			log.Println("EOF", err)
			panic("close")
		}

		data := recvBuf[:n]
		codeBuffer.Write(data)

		msg := Msg{}

		if err = dec.Decode(&msg); nil != err {
			log.Printf("failed to decode message; err: %v", err)
			continue
		}

		log.Println("Server send: ", msg)

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