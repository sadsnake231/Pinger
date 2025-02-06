package main

import (
    "github.com/gin-gonic/gin"
)
func main() {
    router := gin.New()
    router.Use(gin.Logger())
    router.Use(func(c *gin.Context){
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
        c.Next()
    })

    router.POST("/", UpdatePings())
    router.GET("/", GetPings())

    router.Run(":5000")
}

