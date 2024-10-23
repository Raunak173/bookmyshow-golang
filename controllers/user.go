package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/raunak173/bms-go/helpers"
	"github.com/raunak173/bms-go/initializers"
	"github.com/raunak173/bms-go/models"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

type UserRequestBody struct {
	Name        string `json:"name" gorm:"not null" validate:"required,min=2,max=50"`
	Email       string `json:"email" gorm:"not null;unique" validate:"email,required"`
	Password    string `json:"password" gorm:"not null"  validate:"required"`
	PhoneNumber string `json:"phone_number" gorm:"not null" validate:"required"`
}

func SignUp(c *gin.Context) {
	var body UserRequestBody

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := validate.Struct(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash the password"})
		return
	}

	// Create the user
	user := models.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: string(hash),
		IsAdmin:  false,
	}

	// Save the user to the database
	result := initializers.Db.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create the user"})
		return
	}

	// Generate and send OTP
	if _, err := helpers.SendOtp(body.PhoneNumber); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}

	// Return the created user as a response
	c.JSON(http.StatusCreated, gin.H{
		"user":    user,
		"message": "OTP sent to your phone. Please verify.",
	})
}

func VerifyOTP(c *gin.Context) {
	var body struct {
		PhoneNumber string `json:"phone_number" validate:"required"`
		Otp         string `json:"otp" validate:"required,len=6"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := validate.Struct(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// Check OTP using Twilio Verify
	err := helpers.CheckOtp(body.PhoneNumber, body.Otp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		return
	}

	// OTP verified successfully
	c.JSON(http.StatusOK, gin.H{"message": "OTP verification successful"})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=5"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := validate.Struct(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	var user models.User
	initializers.Db.Where("email = ?", body.Email).First(&user)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,                                   // Subject (User ID)
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // Token expiration time
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	// Set the token in a cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*7, "", "", false, true)

	// Send OTP for login
	_, otpErr := helpers.SendOtp(user.PhoneNumber)
	if otpErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent for verification",
		"user":    user,
	})
}
