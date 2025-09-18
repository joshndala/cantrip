package services

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
	Timestamp   string                 `json:"timestamp"`
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

// ProcessChatMessageStream processes a user message with streaming response
func ProcessChatMessageStream(message string, session *ConversationSession, c *gin.Context) error {
	// Call streaming LangGraph agent
	err := callLangGraphAgentStream(message, session, c)
	if err != nil {
		return fmt.Errorf("failed to process message with AI agent: %w", err)
	}

	return nil
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
	// Prepare the request data
	requestData := map[string]interface{}{
		"message":    message,
		"session_id": session.SessionID,
		"context":    session.Context,
		"history":    session.History,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make HTTP request to LangGraph agent
	// In Docker: use service name, in development: use localhost
	agentURL := "http://cantrip-agent:8001/chat"
	if os.Getenv("DOCKER_ENV") == "" {
		agentURL = "http://localhost:8001/chat"
	}
	resp, err := http.Post(agentURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		// Log the error for debugging
		fmt.Printf("Error calling LangGraph agent: %v\n", err)
		// For now, return a mock response since LangGraph agent might not be running
		return &ChatResponse{
			Response:    "I'm your AI Canadian travel assistant! I can help you plan trips across Canada, suggest destinations, create itineraries, and more. What would you like to know?",
			SessionID:   session.SessionID,
			Intent:      "greeting",
			Confidence:  0.8,
			Suggestions: []string{"Tell me about popular Canadian destinations", "Help me plan a trip to Toronto", "What's the weather like in Vancouver?"},
			Timestamp:   time.Now().Format(time.RFC3339),
		}, nil
	}
	defer resp.Body.Close()

	// Parse response
	var response ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Printf("Error decoding LangGraph response: %v\n", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Printf("Successfully got response from LangGraph agent: %s\n", response.Response[:50])
	return &response, nil
}

// callLangGraphAgentStream calls the LangGraph agent for streaming message processing
func callLangGraphAgentStream(message string, session *ConversationSession, c *gin.Context) error {
	// Prepare the request data
	requestData := map[string]interface{}{
		"message":    message,
		"session_id": session.SessionID,
		"context":    session.Context,
		"history":    session.History,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make HTTP request to streaming LangGraph agent
	agentURL := "http://cantrip-agent:8001/chat/stream"
	if os.Getenv("DOCKER_ENV") == "" {
		agentURL = "http://localhost:8001/chat/stream"
	}

	// Create HTTP client with proper timeout for streaming
	client := &http.Client{
		Timeout: 0, // No timeout for streaming
	}

	req, err := http.NewRequest("POST", agentURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error calling streaming LangGraph agent: %v\n", err)
		// Send fallback response
		fallbackResponse := "I'm your AI Canadian travel assistant! I can help you plan trips across Canada, suggest destinations, create itineraries, and more. What would you like to know?"
		words := strings.Fields(fallbackResponse)
		for _, word := range words {
			chunk := map[string]interface{}{
				"type":       "token",
				"content":    word + " ",
				"session_id": session.SessionID,
				"intent":     "greeting",
				"confidence": 0.8,
				"timestamp":  time.Now().Format(time.RFC3339),
			}
			chunkJSON, _ := json.Marshal(chunk)
			fmt.Fprintf(c.Writer, "data: %s\n\n", chunkJSON)
			c.Writer.Flush()
			time.Sleep(50 * time.Millisecond)
		}

		// Send done signal
		doneChunk := map[string]interface{}{
			"type":       "done",
			"session_id": session.SessionID,
		}
		doneJSON, _ := json.Marshal(doneChunk)
		fmt.Fprintf(c.Writer, "data: %s\n\n", doneJSON)
		c.Writer.Flush()
		return nil
	}
	defer resp.Body.Close()

	// Stream the response
	scanner := bufio.NewScanner(resp.Body)
	var fullResponse strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data == "" {
				continue
			}

			// Parse the chunk
			var chunk map[string]interface{}
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue
			}

			// Forward the chunk to the client
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			c.Writer.Flush()

			// Collect the full response for session update
			if chunkType, ok := chunk["type"].(string); ok && chunkType == "token" {
				if content, ok := chunk["content"].(string); ok {
					fullResponse.WriteString(content)
				}
			}

			// If this is the done signal, update the session
			if chunkType, ok := chunk["type"].(string); ok && chunkType == "done" {
				// Update session with the complete response
				UpdateSession(session.SessionID, message, fullResponse.String())
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stream: %w", err)
	}

	return nil
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
