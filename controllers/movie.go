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

func GetMovieByID(c *gin.Context) {
	movieID := c.Param("id")
	// Declare a variable to hold the movie data
	var movie models.Movie
	// Retrieve the movie with its associated venues using GORM's Preload
	if err := initializers.Db.Preload("Venues").First(&movie, movieID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}
	// Return the movie along with its venues
	c.JSON(http.StatusOK, gin.H{
		"movie": movie,
	})
}

type VenueWithShowTimes struct {
	VenueID   uint              `json:"venue_id"`
	VenueName string            `json:"venue_name"`
	ShowTimes []models.ShowTime `json:"show_times"`
}

func GetVenuesByMovieID(c *gin.Context) {
	movieID := c.Param("id")
	var showTimes []models.ShowTime
	/// Retrieve the show times for the given movie ID, preloading the associated venue and movie
	if err := initializers.Db.Preload("Venue").Preload("Movie").Where("movie_id = ?", movieID).Find(&showTimes).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No showtimes found for this movie"})
		return
	}
	// Create a map to club venues with their show times
	venueMap := make(map[uint]gin.H)
	for _, showTime := range showTimes {
		venueID := showTime.Venue.ID
		// Check if the venue is already added to the map
		if venue, exists := venueMap[venueID]; exists {
			// If venue already exists, append the new show time
			venue["show_times"] = append(venue["show_times"].([]string), showTime.Timing)
		} else {
			// If venue doesn't exist, add it to the map with the first show time
			venueMap[venueID] = gin.H{
				"id":         showTime.Venue.ID,
				"name":       showTime.Venue.Name,
				"location":   showTime.Venue.Location,
				"show_times": []string{showTime.Timing},
			}
		}
	}
	// Convert the map to a list for the response
	var venues []gin.H
	for _, venue := range venueMap {
		venues = append(venues, venue)
	}
	c.JSON(http.StatusOK, gin.H{
		"movie_id": movieID,
		"venues":   venues,
	})
}
