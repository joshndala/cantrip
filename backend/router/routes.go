package router

import (
	"github.com/gin-gonic/gin"
	"github.com/joshndala/cantrip/handlers"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *gin.Engine) {
	// API v1 group
	v1 := r.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "healthy"})
		})

		// Chat routes
		chat := v1.Group("/chat")
		{
			chat.POST("/", handlers.ChatHandler)
			chat.GET("/history/:session_id", handlers.GetConversationHistory)
			chat.DELETE("/history/:session_id", handlers.ClearConversation)
			chat.GET("/suggestions/:session_id", handlers.GetConversationSuggestions)
		}

		// Explore routes
		explore := v1.Group("/explore")
		{
			explore.POST("/", handlers.ExploreHandler)
			explore.GET("/mood/:mood", handlers.GetExploreByMood)
		}

		// Itinerary routes
		itinerary := v1.Group("/itinerary")
		{
			itinerary.POST("/", handlers.CreateItineraryHandler)
			itinerary.GET("/:id", handlers.GetItineraryHandler)
			itinerary.PUT("/:id", handlers.UpdateItineraryHandler)
			itinerary.DELETE("/:id", handlers.DeleteItineraryHandler)
		}

		// Packing routes
		packing := v1.Group("/packing")
		{
			packing.POST("/", handlers.GeneratePackingListHandler)
			packing.GET("/:id", handlers.GetPackingListHandler)
			packing.PUT("/:id", handlers.UpdatePackingListHandler)
			packing.GET("/suggestions", handlers.GetPackingSuggestionsHandler)
			packing.GET("/:id/export", handlers.ExportPackingListHandler)
		}

		// Tips routes
		tips := v1.Group("/tips")
		{
			tips.POST("/", handlers.GetTravelTipsHandler)
			tips.GET("/cultural/:destination", handlers.GetCulturalTipsHandler)
			tips.GET("/tipping/:destination", handlers.GetTippingGuideHandler)
			tips.GET("/safety/:destination", handlers.GetSafetyTipsHandler)
			tips.GET("/customs/:destination", handlers.GetLocalCustomsHandler)
			tips.GET("/emergency/:destination", handlers.GetEmergencyInfoHandler)
			tips.GET("/language/:destination", handlers.GetLanguageInfoHandler)
		}

		// Weather routes
		weather := v1.Group("/weather")
		{
			weather.GET("/current", handlers.GetWeatherHandler)
			weather.GET("/forecast", handlers.GetWeatherForecastHandler)
			weather.GET("/forecast/with-notes", handlers.GetWeatherForecastWithNotesHandler)
		}

		// Places routes
		places := v1.Group("/places")
		{
			places.GET("/events", handlers.GetEventsHandler)
			places.GET("/suggestions", handlers.GenerateTripSuggestionsHandler)
		}

		// PDF routes
		pdf := v1.Group("/pdf")
		{
			pdf.POST("/generate", handlers.GeneratePDFHandler)
			pdf.GET("/download/:id", handlers.DownloadPDFHandler)
			pdf.GET("/status/:id", handlers.GetPDFStatusHandler)
			pdf.DELETE("/:id", handlers.DeletePDFHandler)
			pdf.GET("/list", handlers.ListPDFsHandler)
			pdf.POST("/share/:id", handlers.SharePDFHandler)
		}
	}

	// Root route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to CanTrip API",
			"version": "1.0.0",
			"endpoints": gin.H{
				"chat":      "/api/v1/chat",
				"explore":   "/api/v1/explore",
				"itinerary": "/api/v1/itinerary",
				"packing":   "/api/v1/packing",
				"tips":      "/api/v1/tips",
				"weather":   "/api/v1/weather",
				"places":    "/api/v1/places",
				"pdf":       "/api/v1/pdf",
			},
		})
	})
}
