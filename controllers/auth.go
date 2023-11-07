package controllers

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type Auth struct {
	Conf                   *oauth2.Config
	Discord_callback_state string
}

func (a *Auth) Login(c *gin.Context) {
	log.Println(a.Conf.RedirectURL)
	log.Println(a.Conf.ClientID)
	log.Println(a.Conf.ClientSecret)
	log.Println(a.Conf.Scopes)
	log.Println(a.Conf.Endpoint)
	log.Println(os.Getenv("DISCORD_APP_ID"))
	c.Redirect(http.StatusTemporaryRedirect, a.Conf.AuthCodeURL(a.Discord_callback_state))
}

func (a *Auth) DiscordCallback(c *gin.Context) {
	// ensure returned state matches expected state
	if c.PostForm("state") != a.Discord_callback_state {
		c.JSON(http.StatusBadRequest, []byte("State does not match."))
	}

	// grab access token
	token, err := a.Conf.Exchange(context.Background(), c.PostForm("code"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, []byte(err.Error()))
	}

	// access user from token
	res, err := a.Conf.Client(context.Background(), token).Get("https://discord.com/api/users/@me")
	if err != nil || res.StatusCode != 200 {
		if err != nil {
			c.JSON(http.StatusInternalServerError, []byte(err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, res.Status)
		}
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, []byte(err.Error()))
	}
	c.JSON(http.StatusOK, gin.H{"data": body})
}
