package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joshndala/cantrip/services"
)

type TipsRequest struct {
	Destination string   `json:"destination" binding:"required"`
	Category    string   `json:"category"` // "cultural", "practical", "safety", "food", "transport"
	Topics      []string `json:"topics"`
	Language    string   `json:"language"`
}

type TipsResponse struct {
	Destination string             `json:"destination"`
	Category    string             `json:"category"`
	Tips        []services.Tip     `json:"tips"`
	Emergency   services.Emergency `json:"emergency"`
	Language    services.Language  `json:"language"`
}

type Phrase struct {
	English       string `json:"english"`
	Local         string `json:"local"`
	Pronunciation string `json:"pronunciation"`
	Category      string `json:"category"`
}

// GetTravelTipsHandler returns cultural and practical travel tips
func GetTravelTipsHandler(c *gin.Context) {
	var req TipsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get travel tips for the destination
	tips, err := services.GetTravelTips(req.Destination, req.Category, req.Topics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get travel tips"})
		return
	}

	// Get emergency information
	emergency, err := services.GetEmergencyInfo(req.Destination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get emergency information"})
		return
	}

	// Get language information
	language, err := services.GetLanguageInfo(req.Destination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get language information"})
		return
	}

	response := TipsResponse{
		Destination: req.Destination,
		Category:    req.Category,
		Tips:        tips,
		Emergency:   emergency,
		Language:    language,
	}

	c.JSON(http.StatusOK, response)
}

// GetCulturalTipsHandler returns cultural-specific tips
func GetCulturalTipsHandler(c *gin.Context) {
	destination := c.Param("destination")
	if destination == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Destination parameter is required"})
		return
	}

	tips, err := services.GetCulturalTips(destination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cultural tips"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"destination": destination,
		"tips":        tips,
	})
}

// GetTippingGuideHandler returns tipping customs for a destination
func GetTippingGuideHandler(c *gin.Context) {
	destination := c.Param("destination")
	if destination == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Destination parameter is required"})
		return
	}

	tippingGuide, err := services.GetTippingGuide(destination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tipping guide"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"destination":   destination,
		"tipping_guide": tippingGuide,
	})
}

// GetSafetyTipsHandler returns safety tips for a destination
func GetSafetyTipsHandler(c *gin.Context) {
	destination := c.Param("destination")
	if destination == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Destination parameter is required"})
		return
	}

	safetyTips, err := services.GetSafetyTips(destination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get safety tips"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"destination": destination,
		"safety_tips": safetyTips,
	})
}

// GetLocalCustomsHandler returns local customs and etiquette
func GetLocalCustomsHandler(c *gin.Context) {
	destination := c.Param("destination")
	if destination == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Destination parameter is required"})
		return
	}

	customs, err := services.GetLocalCustoms(destination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get local customs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"destination": destination,
		"customs":     customs,
	})
}

// GetEmergencyInfoHandler returns emergency information for a destination
func GetEmergencyInfoHandler(c *gin.Context) {
	destination := c.Param("destination")
	if destination == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Destination parameter is required"})
		return
	}

	emergency, err := services.GetEmergencyInfo(destination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get emergency info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"destination": destination,
		"emergency":   emergency,
	})
}

// GetLanguageInfoHandler returns language information for a destination
func GetLanguageInfoHandler(c *gin.Context) {
	destination := c.Param("destination")
	if destination == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Destination parameter is required"})
		return
	}

	language, err := services.GetLanguageInfo(destination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get language info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"destination": destination,
		"language":    language,
	})
}
