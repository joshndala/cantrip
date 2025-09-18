package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joshndala/cantrip/services"
)

// GetEventsHandler gets events for a city
func GetEventsHandler(c *gin.Context) {
	city := c.Query("city")
	mood := c.Query("mood")
	interests := c.QueryArray("interests")
	date := c.Query("date")

	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city parameter is required"})
		return
	}

	events, err := services.GetEvents(city, mood, interests)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get events: " + err.Error()})
		return
	}

	// Filter events by date if provided
	if date != "" {
		events = services.FilterEventsByDate(events, date)
	}

	// Ensure we return an empty array instead of null
	if events == nil {
		events = []services.Event{}
	}

	c.JSON(http.StatusOK, events)
}

// GenerateTripSuggestionsHandler generates trip suggestions for a city
func GenerateTripSuggestionsHandler(c *gin.Context) {
	city := c.Query("city")
	mood := c.Query("mood")
	interests := c.QueryArray("interests")
	budgetStr := c.Query("budget")
	durationStr := c.Query("duration")

	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city parameter is required"})
		return
	}

	// Parse budget and duration if provided
	var budget float64
	var duration int
	var err error

	if budgetStr != "" {
		budget, err = strconv.ParseFloat(budgetStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid budget parameter"})
			return
		}
	}

	if durationStr != "" {
		duration, err = strconv.Atoi(durationStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid duration parameter"})
			return
		}
	}

	// Get weather info for the city to pass to trip suggestions
	weather, err := services.GetWeather(city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get weather for trip suggestions: " + err.Error()})
		return
	}

	suggestions, err := services.GenerateTripSuggestions(mood, city, budget, duration, interests, weather)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate trip suggestions: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, suggestions)
}
