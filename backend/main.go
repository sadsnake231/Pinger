package main

import (
    "github.com/gin-gonic/gin"
)
func main() {
    router := gin.New()
    router.Use(gin.Logger())

    router.POST("/", UpdatePings())

    /*router.POST("/pings", API.PostPings())
        */

    router.Run(":3000")
}