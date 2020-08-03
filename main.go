package main

import (
	//https://mingrammer.com/translation-go-walkthrough-encoding-package/
	//https://bitlog.tistory.com/124
	"fmt"
	"io"
	"log"
	"net"
)

//type MsgBody struct {
//	Content string
//}
//
//type Msg struct {
//	Header MsgHeader
//	Body   interface{}
//}
//type MsgHeader struct {
//	MsgType string
//	Date    string
//}

type Client struct{
	conn net.Conn
	server *Server
	send chan []byte
}

type Server struct {
	clientMap map[*Client]bool
	ChanEnter chan *Client
	ChanLeave chan *Client
}


func initServer () *Server {
	return &Server{
		clientMap: make(map[*Client]bool),
	}
}

func (s *Server) run(l net.Listener){
	s.ChanEnter = make(chan *Client)
	s.ChanLeave = make(chan *Client)

	for  {
		conn, err := l.Accept()
		if nil != err {
			log.Println(err)
			continue
		}
		client := Client{
			conn: conn,
			server: s,
			send: make(chan []byte),
		}
		s.clientMap[&client] = true
		defer conn.Close()

		go ConnHandler(client)
	}


}


func init() {
	//gob.Register(MsgBody{})
}

func main(){
	fmt.Println("chat start")
	l, err := net.Listen("tcp", ":8000")
	if nil != err {
		log.Println(err)
	}

	defer l.Close()
	server := initServer()
	go server.run(l)

	select{}
}


func ConnHandler(c Client){
	//var codeBuffer bytes.Buffer
	//var dec        *gob.Decoder = gob.NewDecoder(&codeBuffer)
	recvBuf := make([]byte, 4096)

	for {
		n, err := c.conn.Read(recvBuf)
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
			//codeBuffer.Write(data)

			//msg := Msg{}
			//
			//if err = dec.Decode(&msg); nil != err {
			//	log.Printf("failed to decode message; err: %v", err)
			//	continue
			//}

			log.Println("msg: ", string(data))

			// broadcast
			for client := range c.server.clientMap {
				//fmt.Println("client", client)
				//fmt.Println("data", data)
				//fmt.Println("=====================")
				_, err = client.conn.Write(data)
				if err != nil {
					log.Println(err)
					return
				}
			}
			//codeBuffer.Reset()
		}
	}
}