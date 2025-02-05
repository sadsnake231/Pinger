package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.New()
	router.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "pong"})
    })
	router.Run(":8081")
}