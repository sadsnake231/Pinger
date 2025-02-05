package main

func main() {
    router := gin.New()
    router.Use(gin.Logger())

    router.GET("/", func(c *gin.Context()) {
        API.GetPings()
    })

    router.POST("/pings", func(c *gin.Context()) {
        API.PostPings()
    })
}