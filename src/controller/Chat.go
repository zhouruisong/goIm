package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Client struct {
	Uuid string
	Conn *websocket.Conn
}

type ConnManger struct {
	clients     map[string]*Client
	resgister   chan *Client
	unresgister chan *Client
	message     chan []byte
}

type Message struct {
	MessageContext string `json:"messageContext"`
	FromUserUUid   string `json:"fromUserUuid"`
	ToUser         string `json:"toUser"`
	FromUserName   string `json:"fromUserName"`
}

var Manager = ConnManger{
	clients:     make(map[string]*Client),
	resgister:   make(chan *Client),
	unresgister: make(chan *Client),
	message:     make(chan []byte),
}

func (manger *ConnManger) WebSocketStart() {
	for {
		select {
		case conn := <-manger.resgister:
			//有新连接进来后保存
			manger.clients[conn.Uuid] = conn
		case conn := <-manger.unresgister:
			//关闭连接
			conn.Conn.Close()
			delete(manger.clients, conn.Uuid)
			fmt.Println("离开了：", conn.Uuid)
		case mes := <-manger.message:
			Manager.Send(mes)
		}
	}
}

func (manger *ConnManger) Send(mes []byte) {
	messageInfo := manger.ParseMessage(mes)
	//大厅消息
	if messageInfo.ToUser == "all" {
		for _, conn := range manger.clients {
			//本人不广播
			if err := conn.Conn.WriteMessage(1, []byte(messageInfo.MessageContext)); err != nil {
				log.Fatalf("发送消息遇到了错误:%v", err)
			}

			//if conn.Uuid != messageInfo.FromUserUUid {
			//	err := conn.Conn.WriteMessage(1, []byte(messageInfo.MessageContext))
			//	fmt.Println(err)
			//}
		}
	} else if toConn, ok := manger.clients[messageInfo.FromUserUUid]; ok { //私信消息
		toConn.Conn.WriteMessage(1, []byte(messageInfo.MessageContext))

	} else {
		fmt.Println("未知消息")
	}

	//msgString := string(mes)
}

func (manger *ConnManger) ParseMessage(message []byte) Message {
	mes := Message{}
	json.Unmarshal(message, &mes)
	return mes
}

func Ws(c *gin.Context) {
	res := c.Writer
	req := c.Request

	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		//http.NotFound(res, req)
		fmt.Println(error)
		return
	}

	sess, _ := GlobalSessions.SessionStart(c.Writer, c.Request)
	defer sess.SessionRelease(c.Writer)
	info := sess.Get("uuid")
	uid, ok := info.(string)
	if !ok {

		conn.Close()
	}

	Manager.resgister <- &Client{Uuid: uid, Conn: conn}

	go read(conn)

}

func read(conn *websocket.Conn) {

	defer func() {
		conn.Close()
	}()
	for {
		_, s, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			return
		}
		Manager.message <- s
	}

}
