package main

import (
	"os"

	"github.com/djwolff/matchmaker/config"
	"github.com/djwolff/matchmaker/controllers"
	"github.com/djwolff/matchmaker/discord"
	"github.com/djwolff/matchmaker/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func main() {
	r := gin.Default()
	config.Setup()
	models.ConnectDatabase()

	dAuth := controllers.Auth{
		Conf: &oauth2.Config{
			RedirectURL:  os.Getenv("APP_URL") + "/discord/auth/callback",
			ClientID:     os.Getenv("DISCORD_APP_ID"),
			ClientSecret: os.Getenv("DISCORD_APP_SECRET"),
			Scopes:       []string{discord.ScopeIdentify},
			Endpoint:     discord.Endpoint,
		},
		Discord_callback_state: "random",
	}
	r.GET("/login", dAuth.Login)
	r.POST("/discord/auth/callback", dAuth.DiscordCallback)
	r.GET("/users", controllers.FindUsers)

	r.Run()
}
