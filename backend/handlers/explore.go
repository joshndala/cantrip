package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joshndala/cantrip/services"
)

type ExploreRequest struct {
	Mood      string   `json:"mood" binding:"required"`
	City      string   `json:"city" binding:"required"`
	Budget    float64  `json:"budget"`
	Duration  int      `json:"duration"` // in days
	Interests []string `json:"interests"`
	Season    string   `json:"season"`
}

type ExploreResponse struct {
	Suggestions []services.TripSuggestion `json:"suggestions"`
	Weather     services.WeatherInfo      `json:"weather"`
	Events      []services.Event          `json:"events"`
}

// ExploreHandler handles mood and place-based trip suggestions
func ExploreHandler(c *gin.Context) {
	var req ExploreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get weather information
	weather, err := services.GetWeather(req.City)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get weather data"})
		return
	}

	// Get events and attractions
	events, err := services.GetEvents(req.City, req.Mood, req.Interests)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get events data"})
		return
	}

	// Generate trip suggestions based on mood and interests
	suggestions, err := services.GenerateTripSuggestions(req.Mood, req.City, req.Budget, req.Duration, req.Interests, weather)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate suggestions"})
		return
	}

	response := ExploreResponse{
		Suggestions: suggestions,
		Weather:     weather,
		Events:      events,
	}

	c.JSON(http.StatusOK, response)
}

// GetExploreByMood returns suggestions for a specific mood
func GetExploreByMood(c *gin.Context) {
	mood := c.Param("mood")
	city := c.Query("city")

	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "City parameter is required"})
		return
	}

	// Get cached suggestions or generate new ones
	suggestions, err := services.GetCachedSuggestions(mood, city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get suggestions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mood":        mood,
		"city":        city,
		"suggestions": suggestions,
	})
}
