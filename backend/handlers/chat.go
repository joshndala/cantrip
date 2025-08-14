package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joshndala/cantrip/services"
)

type ChatRequest struct {
	Message   string `json:"message" binding:"required"`
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id,omitempty"`
}

// ChatHandler handles conversational interactions
func ChatHandler(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get or create conversation session
	session, err := services.GetOrCreateSession(req.SessionID, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to manage session"})
		return
	}

	// Process message with AI agent
	response, err := services.ProcessChatMessage(req.Message, session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process message"})
		return
	}

	// Update session with new message
	err = services.UpdateSession(session.SessionID, req.Message, response.Response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetConversationHistory returns the conversation history for a session
func GetConversationHistory(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	history, err := services.GetConversationHistory(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get conversation history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id": sessionID,
		"history":    history,
	})
}

// ClearConversation clears the conversation history
func ClearConversation(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	err := services.ClearConversation(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear conversation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Conversation cleared successfully"})
}

// GetConversationSuggestions returns suggested follow-up questions
func GetConversationSuggestions(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	suggestions, err := services.GetConversationSuggestions(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get suggestions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id":  sessionID,
		"suggestions": suggestions,
	})
}
