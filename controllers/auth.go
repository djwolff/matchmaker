package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type Auth struct {
	Conf                   *oauth2.Config
	Discord_callback_state string
}

func (a *Auth) Login(c *gin.Context) {
	oauthState := generateStateOauthCookie(c.Writer)
	c.Redirect(http.StatusTemporaryRedirect, a.Conf.AuthCodeURL(oauthState))
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func (a *Auth) DiscordCallback(c *gin.Context) {
	oauthState, _ := c.Cookie("oauthstate")

	// ensure returned state matches expected state
	if c.Request.FormValue("state") != oauthState {
		fmt.Println("state does not match")
		c.JSON(http.StatusBadRequest, []byte("State does not match."))
		return
	}

	// grab access token
	token, err := a.Conf.Exchange(context.Background(), c.Request.FormValue("code"))
	if err != nil || token == nil {
		if err != nil {
			c.JSON(http.StatusInternalServerError, []byte(err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, []byte("Nil token"))
		}
		return
	}

	// access user from token
	res, err := a.Conf.Client(context.Background(), token).Get("https://discord.com/api/users/@me")
	if err != nil || res == nil || res.StatusCode != 200 {
		if err != nil {
			c.JSON(http.StatusInternalServerError, []byte(err.Error()))
		} else if res == nil {
			c.JSON(http.StatusInternalServerError, []byte("No response callback from discord"))
		} else {
			c.JSON(http.StatusInternalServerError, res.Status)
		}
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, []byte(err.Error()))
		return
	}

	// GetOrCreate User in your db.
	// Redirect or response with a token.
	c.JSON(http.StatusOK, gin.H{"data": body})
}
