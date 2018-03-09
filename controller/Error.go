package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context) {
	c.HTML(http.StatusForbidden, "login.html", nil)
}
