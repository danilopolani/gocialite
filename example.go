package main

import (
	"fmt"
	"net/http"

	"github.com/danilopolani/gocialite"
	"github.com/gin-gonic/gin"
)

var gocial gocialite.Gocial

func main() {
	router := gin.Default()

	router.GET("/", indexHandler)
	router.GET("/auth/:provider", redirectHandler)
	router.GET("/auth/:provider/callback", callbackHandler)

	router.Run("127.0.0.1:9090")
}

// Show homepage with login URL
func indexHandler(c *gin.Context) {
	c.Writer.Write([]byte("<html><head><title>Gocialite example</title></head><body><a href='/auth/github'><button>Login with GitHub</button></a></body></html>"))
}

// Redirect to correct oAuth URL
func redirectHandler(c *gin.Context) {
	// Retrieve provider from route
	provider := c.Param("provider")

	// In this case we use a map to store our secrets, but you can use dotenv or your framework configuration
	// for example, in revel you could use revel.Config.StringDefault(provider + "_clientID", "") etc.
	providerSecrets := map[string]map[string]string{
		"github": {
			"clientID":     "xxxxxxxxxxxxxx",
			"clientSecret": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			"redirectURL":  "http://localhost:9090/auth/github/callback",
		},
		"linkedin": {
			"clientID":     "xxxxxxxxxxxxxx",
			"clientSecret": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			"redirectURL":  "http://localhost:9090/auth/linkedin/callback",
		},
	}

	providerData := providerSecrets[provider]
	authURL := gocial.Driver(provider).
		Scopes([]string{"public_repo"}).
		Redirect(
			providerData["clientID"],
			providerData["clientSecret"],
			providerData["redirectURL"],
		)

	// Redirect with authURL
	c.Redirect(http.StatusFound, authURL)
}

// Handle callback of provider
func callbackHandler(c *gin.Context) {
	// Retrieve provider from route
	provider := c.Param("provider")
	// Retrieve query params for state and code
	state := c.Query("state")
	code := c.Query("code")

	// Handle callback and check for errors
	err := gocial.Handle(state, code)
	if err != nil {
		c.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	// Print in terminal user information
	fmt.Printf("%#v", gocial.User)

	// If no errors, show provider name
	c.Writer.Write([]byte(provider))
}
