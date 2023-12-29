package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/djwolff/matchmaker/db"
	"github.com/djwolff/matchmaker/models"
	"github.com/djwolff/matchmaker/utils/token"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

func Login(c *gin.Context, oauthConfig *oauth2.Config) {
	// Check if the user already has a valid session with an access token
	user, exists := c.Get("user")
	if exists {
		// If the user has a valid session, they are already authenticated
		// You can redirect them to another page or just indicate they are already logged in
		c.JSON(http.StatusOK, gin.H{"message": "User is already logged in", "user": user})
		return
	}

	oauthState := generateStateOauthCookie(c.Writer)

	// Store the generated state in the session
	session := sessions.Default(c)
	session.Set("oauth_state", oauthState)
	session.Save()

	c.Redirect(http.StatusTemporaryRedirect, oauthConfig.AuthCodeURL(oauthState, oauth2.AccessTypeOffline))
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

func DiscordCallback(c *gin.Context, gormDB *gorm.DB, oauthConfig *oauth2.Config) {
	var user models.User
	// Retrieve the stored state from the session
	session := sessions.Default(c)
	storedState, _ := session.Get("oauth_state").(string)

	// Retrieve the state parameter from the callback
	state := c.Query("state")

	// Validate that the retrieved state matches the stored state
	if state != storedState {
		fmt.Println("session state does not match received state")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid state"})
		return
	}

	// Clear the stored state from the session (optional)
	session.Delete("oauth_state")
	session.Save()

	// grab access token
	oauthToken, err := oauthConfig.Exchange(context.Background(), c.Request.FormValue("code"))
	if err != nil || oauthToken == nil {
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		} else {
			c.JSON(http.StatusInternalServerError, "Nil token")
		}
		return
	}

	// access user from token
	res, err := oauthConfig.Client(context.Background(), oauthToken).Get("https://discord.com/api/users/@me")
	if err != nil || res == nil || res.StatusCode != 200 {
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		} else if res == nil {
			c.JSON(http.StatusInternalServerError, "No response callback from discord")
		} else {
			c.JSON(http.StatusInternalServerError, res.Status)
		}
		return
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// save user to database
	savedUser, err := db.GetOrCreateUser(gormDB, user)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// Generate JWT with user claims
	tokenString, err := token.GenerateJWT(savedUser)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// Set the JWT as a cookie
	c.SetCookie("access_token", tokenString, int(time.Hour.Seconds()), "/", "localhost", false, true)

	// TODO: redirect to home page
	c.JSON(http.StatusOK, gin.H{"data": savedUser})
}
