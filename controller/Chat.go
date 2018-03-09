package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gochat/model"

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
		case mes := <-manger.message:
			//消息处理
			Manager.Send(mes)
		}
	}
}

//消息处理
func (manger *ConnManger) Send(mes []byte) {
	messageInfo := manger.ParseMessage(mes)
	//根据uuid  得到 用户信息
	user := model.User{}
	userInfo := user.GetUserByUuid(messageInfo.FromUserUUid)

	messageInfo.FromUserName = userInfo.Nickname
	//大厅消息
	if messageInfo.ToUser == "all" {
		for _, conn := range manger.clients {

			if err := conn.Conn.WriteMessage(1, manger.MessageStructToJson(messageInfo)); err != nil {
				fmt.Printf("发送消息遇到了错误:%v", err)
			}

		}
	} else if toConn, ok := manger.clients[messageInfo.FromUserUUid]; ok { //私信消息
		toConn.Conn.WriteMessage(1, manger.MessageStructToJson(messageInfo))

	} else {
		fmt.Println("未知消息")
	}

	//msgString := string(mes)
}

//解析message 结构体
func (manger *ConnManger) ParseMessage(message []byte) Message {
	mes := Message{}
	json.Unmarshal(message, &mes)
	return mes
}

//结构体转换为json
func (manger *ConnManger) MessageStructToJson(message Message) []byte {
	str, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("结构体转换JSON:%v", err)
	}
	return str
}

//websocket  连接处理
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

	//并发读取
	go read(conn)

}

//读取消息
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
