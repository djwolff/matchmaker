package controllers

import (
	"net/http"

	"github.com/djwolff/matchmaker/models"
	"github.com/gin-gonic/gin"
)

// GET /users
// GET all Users
func FindUsers(c *gin.Context) {
	var users []models.User
	models.DB.Find(&users)

	c.JSON(http.StatusOK, gin.H{"data": users})
}
