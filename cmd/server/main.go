package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Static("/static", "./web")
	r.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})
	r.Run(":8080")
}
