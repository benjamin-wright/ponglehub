package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOriginFunc = func(origin string) bool {
		fmt.Printf("Allowing %s", origin)
		return true
	}
	r.Use(cors.New(config))

	r.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})
	r.Run("0.0.0.0:80")
}
