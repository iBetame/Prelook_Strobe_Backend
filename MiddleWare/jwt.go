package MiddleWare

import (
	"github.com/sunbelife/Prelook_Strobe_Backend/API"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = http.StatusOK
		token := c.GetHeader("token")
		if token == "" {
			code = http.StatusBadRequest
		} else {
			claims, err := API.ParseToken(token)
			if err != nil {
				code = 20001
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = 20002
			}
		}

		if code != http.StatusOK {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  "auth fail",
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}