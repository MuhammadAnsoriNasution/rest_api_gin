package routes

import (
	"fmt"
	"gin-full-rest/config"
	"gin-full-rest/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gopkg.in/danilopolani/gocialite.v1/structs"
)

var JWT_SECRET = "Secret"

// Redirect to correct oAuth URL
func RedirectHandler(c *gin.Context) {
	// Retrieve provider from route
	provider := c.Param("provider")

	providerSecrets := map[string]map[string]string{
		"github": {
			"clientID":     os.Getenv("CLIENT_ID_GH"),
			"clientSecret": os.Getenv("CLIENT_SECRET_GH"),
			"redirectURL":  os.Getenv("AUTH_REDIRECT_URL") + "/github/calback",
		},
	}

	providerScopes := map[string][]string{
		"github": []string{"public_repo"},
	}

	providerData := providerSecrets[provider]
	actualScopes := providerScopes[provider]
	authURL, err := config.Gocial.New().
		Driver(provider).
		Scopes(actualScopes).
		Redirect(
			providerData["clientID"],
			providerData["clientSecret"],
			providerData["redirectURL"],
		)

	// Check for errors (usually driver not valid)
	if err != nil {
		c.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	// Redirect with authURL
	c.Redirect(http.StatusFound, authURL)
}

// Handle callback of provider
func CallbackHandler(c *gin.Context) {
	// Retrieve query params for state and code
	state := c.Query("state")
	code := c.Query("code")
	provider := c.Param("provider")

	// Handle callback and check for errors
	user, token, err := config.Gocial.Handle(state, code)
	if err != nil {
		c.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	var newUser = getOrRegisterUser(provider, (*structs.User)(user))
	var newToken = createToken(newUser)
	c.JSON(200, gin.H{
		"message":  "Berhasil",
		"data":     newUser,
		"token":    token,
		"newtoken": newToken,
	})
	// // Print in terminal user information
	// fmt.Printf("%#v", token)
	// fmt.Printf("%#v", user)

	// // If no errors, show provider name
	// c.Writer.Write([]byte("Hi, " + user.FullName))
}

func createToken(user models.User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"user_role": user.Role,
		"exp":       time.Now().AddDate(0, 0, 7).Unix(),
		"iat":       time.Now().Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		fmt.Println(err)
	}
	return tokenString

}
func getOrRegisterUser(provider string, user *structs.User) models.User {
	var userData models.User
	err := config.DB.Where("provider = ? AND social_id = ? ", provider, user.ID).First(&userData)
	if err.Error != nil {
		newUser := models.User{
			UserName: user.Username,
			FullName: user.FullName,
			Email:    user.Email,
			Provider: provider,
			Avatar:   user.Avatar,
			SocialId: user.ID,
		}
		config.DB.Create(&newUser)
		return newUser
	}
	return userData
}

func CheckToken(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welocme to gin",
	})
}
