package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/djwolff/matchmaker/utils/token"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var (
	// Secret key to sign and verify the JWT
	secretKey = []byte("your-secret-key")

	// OAuth2 configuration
	oauthConfig = &oauth2.Config{
		ClientID:     "your-client-id",
		ClientSecret: "your-client-secret",
		RedirectURL:  "your-redirect-url",
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "authorization-url",
			TokenURL: "token-url",
		},
	}

	// Create a new JWT middleware
	jwtMiddleware = JwtMiddleware(secretKey)
)

// User represents the user model
type User struct {
	ID       string
	Username string
	Email    string
}

// JwtMiddleware creates a Gin middleware for JWT authentication
func JwtMiddleware(secretKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		user, err := ExtractUserFromToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Attach the user to the context for further processing
		c.Set("user", user)

		c.Next()
	}
}

// GenerateJWT generates a new JWT with user claims
func GenerateJWT(user *User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Hour * 1).Unix(), // Token expires in 1 hour
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ExtractUserFromToken extracts user data from the access token
func ExtractUserFromToken(tokenString string) (*User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method and return the secret key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract user data from claims
		user := &User{
			ID:       claims["id"].(string),
			Username: claims["username"].(string),
			Email:    claims["email"].(string),
		}
		return user, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func main() {
	router := gin.Default()

	// OAuth2 login route
	router.GET("/login", func(c *gin.Context) {
		url := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
		c.Redirect(http.StatusTemporaryRedirect, url)
	})

	// OAuth2 callback route
	router.GET("/callback", func(c *gin.Context) {
		code := c.Query("code")

		token, err := oauthConfig.Exchange(oauth2.NoContext, code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get token"})
			return
		}

		// Extract user information from the OAuth2 token
		user := &User{
			ID: token.Extra("id").(string),
			// Extract other user attributes as needed
		}

		// Generate JWT with user claims
		tokenString, err := GenerateJWT(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate JWT"})
			return
		}

		// Set the JWT as a cookie
		c.SetCookie("access_token", tokenString, int(time.Hour.Seconds()), "/", "localhost", false, true)

		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
	})

	// Protected route using JWT middleware
	router.GET("/protected", jwtMiddleware, func(c *gin.Context) {
		// Access the user from the context
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user.(*token.JWTUser)})
	})

	log.Fatal(router.Run(":8080"))
}
