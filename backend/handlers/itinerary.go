package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joshndala/cantrip/services"
)

type ItineraryRequest struct {
	City          string    `json:"city" binding:"required"`
	StartDate     time.Time `json:"start_date" binding:"required"`
	EndDate       time.Time `json:"end_date" binding:"required"`
	Interests     []string  `json:"interests"`
	Budget        float64   `json:"budget"`
	GroupSize     int       `json:"group_size"`
	Pace          string    `json:"pace"`          // "relaxed", "moderate", "intense"
	Accommodation string    `json:"accommodation"` // "budget", "mid-range", "luxury"
}

type ItineraryResponse struct {
	ID        string    `json:"id"`
	City      string    `json:"city"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Days      []DayPlan `json:"days"`
	TotalCost float64   `json:"total_cost"`
	Summary   string    `json:"summary"`
	CreatedAt time.Time `json:"created_at"`
	PDFURL    string    `json:"pdf_url,omitempty"`
}

type DayPlan struct {
	Day           int         `json:"day"`
	Date          time.Time   `json:"date"`
	Activities    []Activity  `json:"activities"`
	Meals         []Meal      `json:"meals"`
	Accommodation string      `json:"accommodation"`
	Transport     []Transport `json:"transport"`
	Notes         string      `json:"notes"`
}

type Activity struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Cost        float64   `json:"cost"`
	Category    string    `json:"category"`
	BookingURL  string    `json:"booking_url,omitempty"`
}

type Meal struct {
	Type        string    `json:"type"` // breakfast, lunch, dinner, snack
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Time        time.Time `json:"time"`
	Cost        float64   `json:"cost"`
	Cuisine     string    `json:"cuisine"`
	Reservation bool      `json:"reservation"`
}

type Transport struct {
	Type      string    `json:"type"` // walking, public, taxi, rental
	From      string    `json:"from"`
	To        string    `json:"to"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Cost      float64   `json:"cost"`
	Duration  int       `json:"duration"` // in minutes
}

// CreateItineraryHandler generates a complete itinerary using LangGraph agent
func CreateItineraryHandler(c *gin.Context) {
	var req ItineraryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate dates
	if req.StartDate.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start date cannot be in the past"})
		return
	}

	if req.EndDate.Before(req.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End date must be after start date"})
		return
	}

	// Call LangGraph agent to generate itinerary
	itinerary, err := services.GenerateItinerary(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate itinerary: " + err.Error()})
		return
	}

	// Save itinerary to cache/database
	err = services.SaveItinerary(itinerary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save itinerary"})
		return
	}

	c.JSON(http.StatusOK, itinerary)
}

// GetItineraryHandler retrieves a specific itinerary
func GetItineraryHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Itinerary ID is required"})
		return
	}

	itinerary, err := services.GetItinerary(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Itinerary not found"})
		return
	}

	c.JSON(http.StatusOK, itinerary)
}

// UpdateItineraryHandler updates an existing itinerary
func UpdateItineraryHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Itinerary ID is required"})
		return
	}

	var req ItineraryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Regenerate itinerary with updated parameters
	itinerary, err := services.GenerateItinerary(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update itinerary"})
		return
	}

	// Preserve the original ID (this would need proper type assertion in real implementation)
	// itinerary.ID = id

	// Save updated itinerary
	err = services.SaveItinerary(itinerary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save updated itinerary"})
		return
	}

	c.JSON(http.StatusOK, itinerary)
}

// DeleteItineraryHandler deletes an itinerary
func DeleteItineraryHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Itinerary ID is required"})
		return
	}

	err := services.DeleteItinerary(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete itinerary"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Itinerary deleted successfully"})
}
