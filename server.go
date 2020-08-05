package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
)

type ClientInfo struct {
	Room Room
	Name string
}

type Room struct {
	Name string
}

type Msg struct {
	Type string
	Text string
	Info ClientInfo
}

type Client struct{
	conn net.Conn
	server *Server
	sendMsg chan Msg
	ChanEnter chan string
}

type Server struct {
	channelMap map[string]map[*Client]bool
	//ChanLeave chan string
}


func initServer () *Server {
	return &Server{
		channelMap: make(map[string]map[*Client]bool),
	}
}

func (s *Server) run(l net.Listener){
	//s.ChanEnter = make(chan string)
	//s.ChanLeave = make(chan string)

	for  {
		conn, err := l.Accept()
		if nil != err {
			log.Println(err)
			continue
		}
		client := Client{
			conn: conn,
			server: s,
			sendMsg: make(chan Msg),
			ChanEnter: make(chan string),
		}
		//s.channelMap["1"][&client] = true
		defer conn.Close()

		go client.ConnHandler()
		go client.brodcast()
		go client.CreateRoom()
	}
}

func (c *Client) CreateRoom() {
	for {
		select {
		case room := <- c.ChanEnter:
			if v, found := c.server.channelMap[room]; found {
				fmt.Println("join room : ", room, v)
			} else {
				fmt.Println("create room : ", room)
				c.server.channelMap[room] = make(map[*Client]bool)
			}
			c.server.channelMap[room][c] = true

		}
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
			// 언마샬링
			//c.sendMsg <- recvBuf[:n]
			data :=  recvBuf[:n]
			var msg Msg
			json.Unmarshal([]byte(data), &msg)
			log.Println("msg: ", msg)

			// create room or into room
			if msg.Type == "room" {
				c.ChanEnter <- msg.Text
			} else { // message
				c.sendMsg <- msg
			}

		}
	}
}


func (c *Client) brodcast() {
	for {
		select {
		case m := <- c.sendMsg:

			b, _ := json.Marshal(m)

			for client := range c.server.channelMap[m.Info.Room.Name]{
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