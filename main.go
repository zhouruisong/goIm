package main

import (
	"gochat/controller"
	"gochat/model"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {

	db, err := model.InitDb()
	if err != nil {
		log.Fatal("db-err:", err)
	}
	defer db.Close()

	route := gin.New()

	//静态文件
	route.Static("/js", "./static/js")
	route.Static("/css", "./static/css")
	route.Static("/vendor", "./static/ui")

	////中间件过滤
	authorized := route.Group("/")

	authorized.Use(controller.Basic())
	{
		authorized.GET("/tochat", controller.ToChat)
		authorized.GET("/ws", controller.Ws)
	}

	//路由
	route.GET("/login", controller.Login)
	//route.POST("/join", controller.Join)
	route.POST("/dologin", controller.DoLogin)
	//错误
	route.GET("/error", controller.Error)

	chat := controller.Manager
	go chat.WebSocketStart()
	//设置模板目录
	route.LoadHTMLGlob(filepath.Join(getCurrentDirectory(), "view/*"))
	http.ListenAndServe(":8975", route)
}

func getCurrentDirectory() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}
