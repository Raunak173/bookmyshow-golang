package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raunak173/bms-go/initializers"
	"github.com/raunak173/bms-go/models"
)

func GetAllMovies(c *gin.Context) {
	var movies []models.Movie
	limit := 5
	l := c.Query("limit")
	if l != "" {
		parsedLimit, err := strconv.Atoi(l)
		if err == nil {
			limit = parsedLimit
		}
	}
	offset := 0
	o := c.Query("offset")
	if o != "" {
		parsedOffset, err := strconv.Atoi(o)
		if err == nil {
			offset = parsedOffset
		}
	}
	initializers.Db.
		Limit(limit).Offset(offset).Find(&movies)
	var totalMovies int64
	initializers.Db.Model(&models.Movie{}).Count(&totalMovies)
	nextOffset := offset + limit
	if nextOffset >= int(totalMovies) {
		nextOffset = -1 // No more movies to load
	}
	c.JSON(200, gin.H{
		"movies":      movies,
		"next_offset": nextOffset,
	})
}

type MovieRequestBody struct {
	Title       string `json:"title" validate:"required,min=2,max=50"`
	Description string `json:"desc" validate:"required"`
	Duration    string `json:"duration" validate:"required"`
}

func CreateMovie(c *gin.Context) {
	var body MovieRequestBody
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if err := validate.Struct(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	//We are checking we are authorized or not
	user, _ := c.Get("user")
	//We get userDetails, because we need to check that we are admin or not
	userDetails := user.(models.User)
	if !userDetails.IsAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, admin access required"})
		return
	}
	movie := models.Movie{
		Title:       body.Title,
		Description: body.Description,
		Duration:    body.Duration,
	}
	result := initializers.Db.Create(&movie)
	if result.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"movie": movie,
	})
}
