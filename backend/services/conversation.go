package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ChatMessage represents a message in the conversation
type ChatMessage struct {
	Role      string    `json:"role"` // "user" or "assistant"
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// ConversationSession represents a chat session
type ConversationSession struct {
	SessionID   string                 `json:"session_id"`
	Context     map[string]interface{} `json:"context"`
	History     []ChatMessage          `json:"history"`
	CreatedAt   time.Time              `json:"created_at"`
	LastUpdated time.Time              `json:"last_updated"`
}

// ChatResponse represents the AI agent's response
type ChatResponse struct {
	Response    string                 `json:"response"`
	SessionID   string                 `json:"session_id"`
	Intent      string                 `json:"intent"`
	Confidence  float64                `json:"confidence"`
	Suggestions []string               `json:"suggestions,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
}

// In-memory session storage (in production, use Redis or database)
var sessions = make(map[string]*ConversationSession)

// GetOrCreateSession gets an existing session or creates a new one
func GetOrCreateSession(sessionID, userID string) (*ConversationSession, error) {
	if sessionID == "" {
		// Generate a new session ID
		sessionID = fmt.Sprintf("session_%d", time.Now().UnixNano())
	}

	if session, exists := sessions[sessionID]; exists {
		session.LastUpdated = time.Now()
		return session, nil
	}

	// Create new session
	session := &ConversationSession{
		SessionID:   sessionID,
		Context:     make(map[string]interface{}),
		History:     []ChatMessage{},
		CreatedAt:   time.Now(),
		LastUpdated: time.Now(),
	}

	sessions[sessionID] = session
	return session, nil
}

// ProcessChatMessage processes a user message with the AI agent
func ProcessChatMessage(message string, session *ConversationSession) (*ChatResponse, error) {
	// Call LangGraph agent for processing
	response, err := callLangGraphAgent(message, session)
	if err != nil {
		return nil, fmt.Errorf("failed to process message with AI agent: %w", err)
	}

	return response, nil
}

// UpdateSession updates the session with new messages
func UpdateSession(sessionID, userMessage, aiResponse string) error {
	session, exists := sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	// Add user message to history
	session.History = append(session.History, ChatMessage{
		Role:      "user",
		Message:   userMessage,
		Timestamp: time.Now(),
	})

	// Add AI response to history
	session.History = append(session.History, ChatMessage{
		Role:      "assistant",
		Message:   aiResponse,
		Timestamp: time.Now(),
	})

	session.LastUpdated = time.Now()
	return nil
}

// GetConversationHistory returns the conversation history
func GetConversationHistory(sessionID string) ([]ChatMessage, error) {
	session, exists := sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	return session.History, nil
}

// ClearConversation clears the conversation history
func ClearConversation(sessionID string) error {
	session, exists := sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	session.History = []ChatMessage{}
	session.Context = make(map[string]interface{})
	session.LastUpdated = time.Now()

	return nil
}

// GetConversationSuggestions returns suggested follow-up questions
func GetConversationSuggestions(sessionID string) ([]string, error) {
	session, exists := sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	// Call LangGraph agent for suggestions
	suggestions, err := callLangGraphSuggestions(session)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggestions: %w", err)
	}

	return suggestions, nil
}

// callLangGraphAgent calls the LangGraph agent for message processing
func callLangGraphAgent(message string, session *ConversationSession) (*ChatResponse, error) {
	// Make HTTP request to LangGraph agent
	resp, err := http.Post("http://localhost:8000/chat", "application/json", nil)
	if err != nil {
		// For now, return a mock response since LangGraph agent might not be running
		return &ChatResponse{
			Response:    "I'm your AI Canadian travel assistant! I can help you plan trips across Canada, suggest destinations, create itineraries, and more. What would you like to know?",
			SessionID:   session.SessionID,
			Intent:      "greeting",
			Confidence:  0.8,
			Suggestions: []string{"Tell me about popular Canadian destinations", "Help me plan a trip to Toronto", "What's the weather like in Vancouver?"},
			Timestamp:   time.Now(),
		}, nil
	}
	defer resp.Body.Close()

	// Parse response
	var response ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// callLangGraphSuggestions calls the LangGraph agent for suggestions
func callLangGraphSuggestions(session *ConversationSession) ([]string, error) {
	// For now, return default suggestions
	// In production, this would call the LangGraph agent
	return []string{
		"Tell me about popular Canadian destinations",
		"Help me plan a trip to Toronto",
		"What's the weather like in Vancouver?",
		"Create a packing list for my Canadian trip",
		"Give me travel tips for Canada",
	}, nil
}
