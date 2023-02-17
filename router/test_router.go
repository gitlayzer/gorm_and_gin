package router

import "github.com/gin-gonic/gin"

func TestRouter(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	v1.GET("/ping", TestHandler)
}

func TestHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
