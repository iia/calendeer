package main

import (
	"fmt"
	"time"
	"github.com/gin-gonic/gin"
)

func get() {
	tckr := time.NewTicker(4 * time.Second)

	for {
		select {
			case <-tckr.C:
				fmt.Println("GET")
		}
	}
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	go get()
	r.Run(":5000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
