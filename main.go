package main

import (
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/djwolff/matchmaker/config"
	"github.com/djwolff/matchmaker/controllers"
	"github.com/djwolff/matchmaker/db"
	"github.com/djwolff/matchmaker/discord"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func main() {
	config.Setup()
	db := db.ConnectDatabase()
	r := gin.Default()

	OAuthConfig := &oauth2.Config{
		RedirectURL:  os.Getenv("APP_URL") + "/auth/callback",
		ClientID:     os.Getenv("DISCORD_APP_ID"),
		ClientSecret: os.Getenv("DISCORD_APP_SECRET"),
		Scopes:       []string{discord.ScopeIdentify},
		Endpoint:     discord.Endpoint,
	}

	// Use sessions middleware
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// Secure headers
	r.Use(func(c *gin.Context) {
		c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains") // Enforce HTTPS
		c.Header("Content-Security-Policy", "script-src 'self'")                     // Enable Content Security Policy
		c.Next()
	})

	r.GET("/login", func(c *gin.Context) {
		controllers.Login(c, OAuthConfig)
	})
	r.GET("/auth/callback", func(c *gin.Context) {
		controllers.DiscordCallback(c, db, OAuthConfig)
	})
	r.GET("/auth/protected", func(c *gin.Context) {
		controllers.Protected(c)
	})

	r.Run()
}
