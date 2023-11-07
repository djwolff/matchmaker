package main
 
import (
    "log"
    "os"
    "github.com/gin-gonic/gin"
    "github.com/golangcompany/restfulapui/controllers"
)
 
func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port := "8080"
        // }
        router := gin.New()
        router.Use(gin.Logger())
 
        router.POST("/createuser", controllers.CreateUser())
        log.Fatal(router.Run(":" + port))
    }
}