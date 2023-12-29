package token

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/djwolff/matchmaker/models"
	"github.com/gin-gonic/gin"
)

type JWTUser struct {
	ID         string
	Username   string
	GlobalName string
	Avatar     string
}

func GenerateJWT(user *models.User) (string, error) {
	token_lifespan, err := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"id":         user.ID,
		"username":   user.Username,
		"globalname": user.GlobalName,
		"avatar":     user.Avatar,
		"exp":        time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix(),
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func TokenValid(c *gin.Context) error {
	tokenString := ExtractToken(c)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return err
	}
	return nil
}

// TODO: three ways of getting token? kinda sus
func ExtractToken(c *gin.Context) string {
	tokenString, _ := c.Cookie("access_token")
	if tokenString != "" {
		return tokenString
	}
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func ExtractUserFromToken(c *gin.Context) (*JWTUser, error) {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := &JWTUser{
			ID:         claims["id"].(string),
			Username:   claims["username"].(string),
			GlobalName: claims["globalname"].(string),
			Avatar:     claims["avatar"].(string),
		}
		return user, nil
	}
	return nil, fmt.Errorf("invalid token")
}
