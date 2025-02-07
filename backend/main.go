package main

import (
    "log"
    "sync"
    
    "github.com/gin-gonic/gin"
)

func main() {
    conn = EstablishConnection()
    defer conn.Close()

    err := SetupRabbitMQ()
    if err != nil {
        log.Fatalf(err.Error())
    }
    defer CloseRabbitMQ()

    var wg sync.WaitGroup
    wg.Add(1)

    go func(){
        defer wg.Done()
        StartQueueWorker(conn)
    }()


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
    

    wg.Wait()
}