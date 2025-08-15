package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joshndala/cantrip/services"
)

// GetWeatherHandler gets current weather for a city
func GetWeatherHandler(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city parameter is required"})
		return
	}

	weather, err := services.GetWeather(city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get weather: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, weather)
}

// GetWeatherForecastHandler gets weather forecast for a city and date range
func GetWeatherForecastHandler(c *gin.Context) {
	city := c.Query("city")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if city == "" || startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city, start_date, and end_date parameters are required"})
		return
	}

	forecast, err := services.GetWeatherForecast(city, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get weather forecast: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, forecast)
}

// GetWeatherForecastWithNotesHandler gets weather forecast with helpful notes
func GetWeatherForecastWithNotesHandler(c *gin.Context) {
	city := c.Query("city")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if city == "" || startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city, start_date, and end_date parameters are required"})
		return
	}

	forecast, notes, err := services.GetWeatherForecastWithNotes(city, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get weather forecast: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"forecast": forecast,
		"notes":    notes,
	})
}
