package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raunak173/bms-go/helpers"
	"github.com/raunak173/bms-go/initializers"
	"github.com/raunak173/bms-go/models"
)

func GetAllVenues(c *gin.Context) {
	var venues []models.Venue
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
		Limit(limit).Offset(offset).Find(&venues)
	var totalVenues int64
	initializers.Db.Model(&models.Venue{}).Count(&totalVenues)
	nextOffset := offset + limit
	if nextOffset >= int(totalVenues) {
		nextOffset = -1 // No more venues to load
	}
	c.JSON(200, gin.H{
		"venues":      venues,
		"next_offset": nextOffset,
	})
}

type VenueRequestBody struct {
	Name     string `json:"name" validate:"required"`
	Location string `json:"location" validate:"required"`
}

func CreateVenue(c *gin.Context) {
	var body VenueRequestBody
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
	venue := models.Venue{
		Name:     body.Name,
		Location: body.Location,
	}
	result := initializers.Db.Create(&venue)
	if result.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"venue": venue,
	})
}

type MovieIdsBody struct {
	MovieIDs []uint `json:"movie_ids"`
}

func AddMoviesInVenue(c *gin.Context) {
	venueID := c.Param("id")
	var body MovieIdsBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie IDs"})
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
	//Check if venue exists
	var venue models.Venue
	if err := initializers.Db.First(&venue, venueID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Venue not found"})
		return
	}
	// Retrieve all movies by their IDs
	var movies []models.Movie
	if err := initializers.Db.Where("id IN ?", body.MovieIDs).Find(&movies).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Some movies not found"})
		return
	}
	// Associate the movies with the venue
	if err := initializers.Db.Model(&venue).Association("Movies").Append(movies); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movies to venue"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Movies successfully added to venue"})
}

func GetVenueByID(c *gin.Context) {
	venueID := c.Param("id")
	var venue models.Venue
	// Use GORM to preload the Movies relationship while fetching the venue by ID
	if err := initializers.Db.Preload("Movies").First(&venue, venueID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Venue not found"})
		return
	}
	// Respond with the venue and its associated movies
	c.JSON(http.StatusOK, gin.H{
		"venue": venue,
	})
}

type ShowTimingsBody struct {
	ShowTimings []string `json:"show_timings"`
	MovieId     uint     `json:"movie_id"`
}

func AddShowTimings(c *gin.Context) {
	venueID := c.Param("id")
	var body ShowTimingsBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie IDs"})
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
	//Check if venue exists
	var venue models.Venue
	if err := initializers.Db.First(&venue, venueID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Venue not found"})
		return
	}
	// Validate the movie ID
	var movie models.Movie
	if err := initializers.Db.First(&movie, body.MovieId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}
	// Add show timings
	for _, timingStr := range body.ShowTimings {
		// Create a new ShowTime record
		showTime := models.ShowTime{
			Timing:  timingStr,
			MovieID: body.MovieId, // Associate with the movie
			VenueID: venue.ID,     // Associate with the venue
		}
		// Save the show time record
		if err := initializers.Db.Create(&showTime).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error saving show time for %s: %v", timingStr, err)})
			return
		}
		// Generate the default seat layout for this showtime
		seats := helpers.GenerateSeatsForShowTime(showTime.ID)
		// Save the generated seats to the database
		if err := initializers.Db.Create(&seats).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error generating seats for %s: %v", timingStr, err)})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Show timings added successfully"})
}
