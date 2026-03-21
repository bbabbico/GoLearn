package main

import (
	Aidb "awesomeGO/AI/aidb"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	Aidb.Init()

	r.GET("/", func(c *gin.Context) {
		Aidb.Load()
		c.JSON(200, gin.H{})
	})

	r.Run(":8080")
}
