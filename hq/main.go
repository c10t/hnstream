package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

  r.GET("/start", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "start reading",
		})
	})

  r.GET("/stop", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "stop reading",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
