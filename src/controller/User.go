package controller

import (
	"gochat/src/model"
	"net/http"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

func DoLogin(c *gin.Context) {
	name := c.PostForm("nickname")
	password := c.PostForm("password")
	//解析中文长度
	if nameLen := utf8.RuneCount([]byte(name)); nameLen >= 16 {
		c.JSON(http.StatusAccepted, gin.H{"result": false, "message": "昵称不能超多16个字符的"})
		return
	}

	if passwordLen := utf8.RuneCount([]byte(password)); passwordLen < 6 {
		c.JSON(http.StatusAccepted, gin.H{"result": false, "message": "密码六位数的"})
		return
	}

	user := model.User{}
	userInfo := user.GetUserByName(name)
	//密码相同说明是同一账号
	if userInfo.Password != password {

		if ok := user.CreateUser(name, password); !ok {
			c.JSON(http.StatusAccepted, gin.H{"result": false, "message": "登入失败"})
		}
		userInfo = user.GetUserByName(name)
	}

	if saveSession(c, userInfo.Uuid) {
		c.JSON(http.StatusOK, gin.H{
			"result": true,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"result": false,
		})
	}
}

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func saveSession(c *gin.Context, uuid string) bool {
	sess, _ := GlobalSessions.SessionStart(c.Writer, c.Request)
	defer sess.SessionRelease(c.Writer)

	err := sess.Set("uuid", uuid)
	if err != nil {
		return false
	}

	return true
}
