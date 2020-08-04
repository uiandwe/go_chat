package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
)

type Msg struct {
	Type string
	Text string
}

type Client struct{
	conn net.Conn
	server *Server
	send chan []byte
}

type Server struct {
	channelMap map[string]map[*Client]bool
	ChanEnter chan *Client
	ChanLeave chan *Client
}


func initServer () *Server {
	return &Server{
		channelMap: make(map[string]map[*Client]bool),
	}
}

func (s *Server) run(l net.Listener){
	s.ChanEnter = make(chan *Client)
	s.ChanLeave = make(chan *Client)
	s.channelMap["1"] = make(map[*Client]bool)

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
		s.channelMap["1"][&client] = true
		defer conn.Close()

		go client.ConnHandler()
		go client.brodcast()
	}


}

func (c *Client)ConnHandler(){

	recvBuf := make([]byte, 4096)

	for {
		n, err := c.conn.Read(recvBuf)
		if nil != err {
			if io.EOF == err {

				delete(c.server.channelMap["1"], c)
				log.Println("client close", err)
				return
			}
			log.Println(err)
			return
		}

		if 0 < n {
			c.send <- recvBuf[:n]
		}
	}
}


func (c *Client) brodcast() {
	for {
		select {
		case data := <- c.send:
			var m = Msg {"text", string(data)}
			log.Println("msg: ", m)
			b, _ := json.Marshal(&m)
			for client := range c.server.channelMap["1"]{
				_, err := client.conn.Write(b)
				if err != nil {
					log.Println(err)
				}
			}

		}
	}
}


func init() {
	fmt.Println("init")
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