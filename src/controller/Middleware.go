package controller

import (
	"github.com/gin-gonic/gin"
)

//定义
func Basic() gin.HandlerFunc {
	return func(c *gin.Context) {
		//sess, _ := GlobalSessions.SessionStart(c.Writer, c.Request)
		//defer sess.SessionRelease(c.Writer)
		//if info := sess.Get("user"); info != nil {
		//	c.Next()
		//}
		//c.Redirect(302, "/login")
		//return
	}
}
