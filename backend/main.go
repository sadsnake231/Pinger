package main

import (
    "log"
    "sync"

    "Pinger/api"
    "Pinger/queue"
    "Pinger/database"
    
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
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
    
    router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

    router.POST("/", api.UpdatePings())
    router.GET("/", api.AuthMiddleware(), api.GetPings())

    router.POST("/signup", api.SignUp())
    router.POST("/login", api.Login())
    router.POST("/logout", api.LogoutUser())

    router.Run(":5000")
    

    wg.Wait()
}