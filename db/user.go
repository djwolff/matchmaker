package db

import (
	"fmt"

	"github.com/djwolff/matchmaker/models"
	"gorm.io/gorm"
)

func GetOrCreateUser(gormDB *gorm.DB, user models.User) (*models.User, error) {
	// if user exists, return user
	// if user not exists, create and return user
	fmt.Println(user.ID)
	// var foundOrCreatedUser models.User
	if err := gormDB.FirstOrCreate(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
