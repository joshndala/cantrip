package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joshndala/cantrip/services"
)

type PackingRequest struct {
	Destination  string   `json:"destination" binding:"required"`
	StartDate    string   `json:"start_date" binding:"required"`
	EndDate      string   `json:"end_date" binding:"required"`
	Activities   []string `json:"activities"`
	Weather      string   `json:"weather"`
	GroupSize    int      `json:"group_size"`
	AgeGroup     string   `json:"age_group"` // "adult", "child", "senior"
	SpecialNeeds []string `json:"special_needs"`
	BaggageType  string   `json:"baggage_type"` // "carry-on", "checked", "both"
}

type PackingResponse struct {
	ID          string               `json:"id"`
	Destination string               `json:"destination"`
	Categories  []PackingCategory    `json:"categories"`
	TotalItems  int                  `json:"total_items"`
	Notes       []string             `json:"notes"`
	Weather     services.WeatherInfo `json:"weather"`
}

type PackingCategory struct {
	Name     string        `json:"name"`
	Items    []PackingItem `json:"items"`
	Priority string        `json:"priority"` // "essential", "recommended", "optional"
	Weight   float64       `json:"weight"`   // in kg
}

type PackingItem struct {
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Essential bool    `json:"essential"`
	Notes     string  `json:"notes"`
	Category  string  `json:"category"`
	Weight    float64 `json:"weight"` // in kg
	Volume    float64 `json:"volume"` // in liters
}

// GeneratePackingListHandler creates a personalized packing list
func GeneratePackingListHandler(c *gin.Context) {
	var req PackingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get weather forecast for the destination
	weather, err := services.GetWeather(req.Destination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get weather data"})
		return
	}

	// Convert to services.PackingRequest
	serviceReq := services.PackingRequest{
		Destination:  req.Destination,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Activities:   req.Activities,
		Weather:      req.Weather,
		GroupSize:    req.GroupSize,
		AgeGroup:     req.AgeGroup,
		SpecialNeeds: req.SpecialNeeds,
		BaggageType:  req.BaggageType,
	}

	// Generate packing list based on destination, weather, and activities
	packingList, err := services.GeneratePackingList(serviceReq, weather)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate packing list"})
		return
	}

	// Save packing list to cache
	err = services.SavePackingList(packingList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save packing list"})
		return
	}

	c.JSON(http.StatusOK, packingList)
}

// GetPackingListHandler retrieves a specific packing list
func GetPackingListHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Packing list ID is required"})
		return
	}

	packingList, err := services.GetPackingList(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Packing list not found"})
		return
	}

	c.JSON(http.StatusOK, packingList)
}

// UpdatePackingListHandler updates an existing packing list
func UpdatePackingListHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Packing list ID is required"})
		return
	}

	var req PackingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get updated weather data
	weather, err := services.GetWeather(req.Destination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get weather data"})
		return
	}

	// Convert to services.PackingRequest
	serviceReq := services.PackingRequest{
		Destination:  req.Destination,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Activities:   req.Activities,
		Weather:      req.Weather,
		GroupSize:    req.GroupSize,
		AgeGroup:     req.AgeGroup,
		SpecialNeeds: req.SpecialNeeds,
		BaggageType:  req.BaggageType,
	}

	// Regenerate packing list
	packingList, err := services.GeneratePackingList(serviceReq, weather)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update packing list"})
		return
	}

	packingList.ID = id // Preserve the original ID

	// Save updated packing list
	err = services.SavePackingList(packingList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save updated packing list"})
		return
	}

	c.JSON(http.StatusOK, packingList)
}

// GetPackingSuggestionsHandler returns packing suggestions for a destination
func GetPackingSuggestionsHandler(c *gin.Context) {
	destination := c.Query("destination")
	season := c.Query("season")
	activities := c.QueryArray("activities")

	if destination == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Destination parameter is required"})
		return
	}

	suggestions, err := services.GetPackingSuggestions(destination, season, activities)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get packing suggestions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"destination": destination,
		"season":      season,
		"activities":  activities,
		"suggestions": suggestions,
	})
}

// ExportPackingListHandler exports packing list as PDF
func ExportPackingListHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Packing list ID is required"})
		return
	}

	packingList, err := services.GetPackingList(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Packing list not found"})
		return
	}

	// Generate PDF
	pdfURL, err := services.GeneratePackingListPDF(packingList.ID, "pdf", true, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pdf_url": pdfURL,
		"message": "Packing list PDF generated successfully",
	})
}
