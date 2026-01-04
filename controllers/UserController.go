package controllers

import (
	"jwt/initilizers"
	"jwt/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	// get the req body
	var body struct {
		Email    string
		Password string
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read the body",
		})
		return
	}

	// hash the password

	hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error is generating the Hashed Password",
		})
		return
	}

	// create the user
	user := models.User{Email: body.Email, Password: string(hashed)}

	result := initilizers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Invalid Username or Password",
		})
		return
	}

	// respond

	c.JSON(http.StatusOK, gin.H{})

}

func Login(c *gin.Context) {
	// get the request body
	var body struct {
		Email    string
		Password string
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to retrive the body",
		})
		return
	}

	// lookup requested user
	var user models.User
	initilizers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Email or Password",
		})
		return
	}

	//compare sent hash pass with the user hash pass
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "The Password is not matched",
		})
		return
	}

	// Generate the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenstring, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to Generate the JWT tokenString",
		})
		return
	}

	// setting the cookies

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenstring, 3600*24*30, "", "", false, true)

	// respond

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})

}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}
