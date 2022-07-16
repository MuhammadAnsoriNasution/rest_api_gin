package midleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func IsAuth() gin.HandlerFunc {
	return checkJwt(false)
}
func IsAdmin() gin.HandlerFunc {
	return checkJwt(true)
}
func checkJwt(middlewareAdmin bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		authheader := c.Request.Header.Get("Authorization")
		bearerToken := strings.Split(authheader, " ")

		if len(bearerToken) == 2 {
			token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
				return []byte(os.Getenv("JWT_SECRET")), nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				// fmt.Println(claims["user_id"], claims["user_role"])
				userRole := claims["user_role"]
				c.Set("jwt_user_id", claims["user_id"])
				// c.Set("jwt_isAdmin", claims["user_role"])
				fmt.Println(userRole)
				fmt.Println(middlewareAdmin)
				if middlewareAdmin && userRole == false {
					c.JSON(403, gin.H{
						"message": "only admin allowed",
					})
					c.Abort()
					return
				}
			} else {
				c.JSON(401, gin.H{
					"message": "gagal",
					"err":     err,
				})
				c.Abort()
				return
			}
		} else {
			c.JSON(401, gin.H{
				"message": "Authorization token not provided",
			})
			c.Abort()
			return
		}
	}

}
