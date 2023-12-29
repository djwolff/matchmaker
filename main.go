package main

import (
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/djwolff/matchmaker/config"
	"github.com/djwolff/matchmaker/controllers"
	"github.com/djwolff/matchmaker/db"
	"github.com/djwolff/matchmaker/discord"
	"github.com/djwolff/matchmaker/middlewares"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func main() {
	config.Setup()
	db := db.ConnectDatabase()
	r := gin.Default()
	matchServer := controllers.NewMatchmakingServer()
	matchServer.ContinuousMatchmaking()

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

	// TODO: homepage
	// r.GET("/", )

	// TODO: video game page
	// r.GET("/", )
	protected := r.Group("/games")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.POST("/:videogame", func(c *gin.Context) {
		matchServer.MatchMake(c, db, c.Param("videogame"))
	})

	// TODO: profile page
	// r.GET("/profile", )

	// TODO: profile page
	// r.GET("/{userID}", )

	// TODO: history
	// r.GET("/history", )

	// TODO: friends page
	// r.GET("/friends", )

	r.Run()
}
