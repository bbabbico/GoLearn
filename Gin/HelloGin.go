package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	p := gin.Default()
	p.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, Gin!",
		})
	})

	p.Run(":8080")
}
