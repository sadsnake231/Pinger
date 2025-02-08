package main

import (
    "log"
    "sync"

    "Pinger/api"
    "Pinger/queue"
    "Pinger/database"
    
    "github.com/gin-gonic/gin"
)

func main() {
    conn := database.EstablishConnection()
    defer conn.Close()

    err := queue.SetupRabbitMQ()
    if err != nil {
        log.Fatalf(err.Error())
    }
    defer queue.CloseRabbitMQ()

    var wg sync.WaitGroup
    wg.Add(1)

    go func(){
        defer wg.Done()
        queue.StartQueueWorker(conn)
    }()


    router := gin.New()
    router.Use(gin.Logger())
    router.Use(func(c *gin.Context){
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
        c.Next()
    })

    router.POST("/", api.UpdatePings())
    router.GET("/", api.AuthMiddleware(), api.GetPings())

    router.POST("/signup", api.SignUp())
    router.POST("/login", api.Login())
    router.POST("/logout", api.LogoutUser())

    router.Run(":5000")
    

    wg.Wait()
}