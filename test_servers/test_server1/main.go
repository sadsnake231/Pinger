package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.New()
	router.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "this is server 1"})
    })
	router.Run(":8081")
}