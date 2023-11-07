package models
 
Imports "go.mongodb.org/mongo-driver/bson/primitive"
 
type User struct {
    Name string `json:"name" bson:"name"`
    Age  int    `json:"age" bson:"age"`
}