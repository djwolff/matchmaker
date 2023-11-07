package controllers
 
import (
    "context"
    "fmt"
    "net/http"
    "time"
 
    "github.com/gin-gonic/gin"
    "github.com/golangcompany/restfulapui/database"
    "github.com/golangcompany/restfulapui/models"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "gopkg.in/mgo.v2/bson"
)
 
var UserCollection *mongo.Collection = database.UserData(database.Client, "User2")
 
func CreateUser() gin.HandlerFunc {
    return func(c *gin.Context) {
 
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()
        var User models.User
 
        if err := c.BindJSON(&User); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            c.Abort()
        }
        User.ID = primitive.NewObjectID()
        _, inserterr := UserCollection.InsertOne(ctx, User)
        if inserterr != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "not created"})
            c.Abort()
        }
        c.IndentedJSON(200, "user created successfully")
    }
}