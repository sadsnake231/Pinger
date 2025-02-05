package main

import (
    "github.com/gin-gonic/gin"
)
func main() {
    router := gin.New()
    router.Use(gin.Logger())

    router.POST("/", UpdatePings())
    router.GET("/", GetPings())

    router.Run(":3000")
}

