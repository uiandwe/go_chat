package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
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

func main(){
	fmt.Println("chat start")
	l, err := net.Listen("tcp", ":8000")
	if nil != err {
		log.Println(err)
	}

	defer l.Close()

	for  {
		conn, err := l.Accept()
		if nil != err {
			log.Println(err)
			continue
		}

		defer conn.Close()
		go ConnHandler(conn)
	}
}

func init() {
	gob.Register(MsgBody{})
}

func ConnHandler(conn net.Conn){
	var codeBuffer bytes.Buffer
	var dec        *gob.Decoder = gob.NewDecoder(&codeBuffer)
	recvBuf := make([]byte, 4096)

	for {
		n, err := conn.Read(recvBuf)
		if nil != err {
			if io.EOF == err {
				log.Println(err)
				return
			}
			log.Println(err)
			return
		}

		if 0 < n {
			data := recvBuf[:n]
			codeBuffer.Write(data)

			msg := Msg{}

			if err = dec.Decode(&msg); nil != err {
				log.Printf("failed to decode message; err: %v", err)
				continue
			}

			log.Println("msg: ", msg)
		}
	}
}