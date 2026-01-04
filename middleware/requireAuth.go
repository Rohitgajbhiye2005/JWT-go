package middleware

import (
	"jwt/initilizers"
	"jwt/models"

	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	// Get the cookie of the request
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
			
		})
		return
	}
	//Decode/validate it

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
		})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	//Check the exp
	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	//Find the user with token sub

	var user models.User
	initilizers.DB.First(&user, claims["sub"])

	if user.ID == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	//Attach to req
	c.Set("user", user)

	//continue
	c.Next()
}
