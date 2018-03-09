package controller

import (
	"net/http"

	"gochat/model"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func Join(c *gin.Context) {
	name := c.PostForm("nickname")
	sess, _ := GlobalSessions.SessionStart(c.Writer, c.Request)
	defer sess.SessionRelease(c.Writer)

	err := sess.Set("user", name)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"result": true,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"result": false,
		})
	}

}

func ToChat(c *gin.Context) {

	sess, _ := GlobalSessions.SessionStart(c.Writer, c.Request)
	defer sess.SessionRelease(c.Writer)

	uuid := sess.Get("uuid")
	user := model.User{}
	if uid, ok := uuid.(string); ok {
		userInfo := user.GetUserByUuid(uid)
		c.HTML(http.StatusOK, "tochat.html", gin.H{
			"uuid":     userInfo.Uuid,
			"username": userInfo.Nickname,
		})
	}

}
