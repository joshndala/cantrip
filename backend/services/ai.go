package services

//COMPLETED
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// AI service configuration
const (
	LangGraphBaseURL = "http://localhost:8001"
	DefaultTimeout   = 30 * time.Second
)

// ItineraryRequest represents a request to generate an itinerary
type ItineraryRequest struct {
	City          string   `json:"city"`
	StartDate     string   `json:"start_date"`
	EndDate       string   `json:"end_date"`
	Interests     []string `json:"interests"`
	Budget        float64  `json:"budget"`
	GroupSize     int      `json:"group_size"`
	Pace          string   `json:"pace"`          // relaxed, moderate, intense
	Accommodation string   `json:"accommodation"` // budget, mid-range, luxury
}

// ItineraryResponse represents the response from itinerary generation
type ItineraryResponse struct {
	Success   bool                   `json:"success"`
	Itinerary map[string]interface{} `json:"itinerary"`
	Metadata  struct {
		City        string  `json:"city"`
		Duration    int     `json:"duration"`
		TotalCost   float64 `json:"total_cost"`
		GeneratedAt string  `json:"generated_at"`
	} `json:"metadata"`
}

// ExploreRequest represents a request to explore a destination
type ExploreRequest struct {
	Mood      string   `json:"mood"`
	City      string   `json:"city"`
	Budget    float64  `json:"budget"`
	Duration  int      `json:"duration"`
	Interests []string `json:"interests"`
}

// ExploreResponse represents the response from destination exploration
type ExploreResponse struct {
	Success     bool                     `json:"success"`
	Suggestions []map[string]interface{} `json:"suggestions"`
	Weather     map[string]interface{}   `json:"weather"`
	Events      []map[string]interface{} `json:"events"`
	Metadata    struct {
		City        string `json:"city"`
		Mood        string `json:"mood"`
		GeneratedAt string `json:"generated_at"`
	} `json:"metadata"`
}

// ChatRequest represents a chat request
type ChatRequest struct {
	Message   string                   `json:"message"`
	SessionID string                   `json:"session_id"`
	Context   map[string]interface{}   `json:"context"`
	History   []map[string]interface{} `json:"history"`
}

// AIChatResponse represents the response from chat
type AIChatResponse struct {
	Response    string                   `json:"response"`
	SessionID   string                   `json:"session_id"`
	Intent      string                   `json:"intent"`
	Confidence  float64                  `json:"confidence"`
	Suggestions []map[string]interface{} `json:"suggestions"`
	Data        map[string]interface{}   `json:"data"`
	Timestamp   string                   `json:"timestamp"`
}

// AIPackingRequest represents a request to generate a packing list
type AIPackingRequest struct {
	Destination  string   `json:"destination"`
	StartDate    string   `json:"start_date"`
	EndDate      string   `json:"end_date"`
	Activities   []string `json:"activities"`
	Weather      string   `json:"weather"`
	GroupSize    int      `json:"group_size"`
	AgeGroup     string   `json:"age_group"`
	SpecialNeeds []string `json:"special_needs"`
}

// AIPackingResponse represents the response from packing list generation
type AIPackingResponse struct {
	Success     bool                   `json:"success"`
	PackingList map[string]interface{} `json:"packing_list"`
	Weather     map[string]interface{} `json:"weather"`
	Metadata    struct {
		Destination string `json:"destination"`
		TotalItems  int    `json:"total_items"`
		GeneratedAt string `json:"generated_at"`
	} `json:"metadata"`
}

// AIClient represents the AI service client
type AIClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewAIClient creates a new AI client
func NewAIClient() *AIClient {
	baseURL := os.Getenv("LANGGRAPH_BASE_URL")
	if baseURL == "" {
		baseURL = LangGraphBaseURL
	}

	return &AIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// Global AI client instance
var aiClient *AIClient

// InitializeAI initializes the global AI client
func InitializeAI() {
	aiClient = NewAIClient()
}

// GetAIClient returns the global AI client
func GetAIClient() *AIClient {
	if aiClient == nil {
		InitializeAI()
	}
	return aiClient
}

// GenerateItinerary generates a complete itinerary using the LangGraph agent
func GenerateItinerary(req ItineraryRequest) (*ItineraryResponse, error) {
	client := GetAIClient()

	// Convert request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal itinerary request: %w", err)
	}

	// Make HTTP request to LangGraph agent
	resp, err := client.makeRequest("POST", "/generate-itinerary", jsonData)
	if err != nil {
		return nil, err
	}

	// Parse response
	var itineraryResp ItineraryResponse
	if err := json.Unmarshal(resp, &itineraryResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal itinerary response: %w", err)
	}

	return &itineraryResp, nil
}

// ExploreDestination explores a destination using the LangGraph agent
func ExploreDestination(req ExploreRequest) (*ExploreResponse, error) {
	client := GetAIClient()

	// Convert request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal explore request: %w", err)
	}

	// Make HTTP request to LangGraph agent
	resp, err := client.makeRequest("POST", "/explore-destination", jsonData)
	if err != nil {
		return nil, err
	}

	// Parse response
	var exploreResp ExploreResponse
	if err := json.Unmarshal(resp, &exploreResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal explore response: %w", err)
	}

	return &exploreResp, nil
}

// Chat handles conversational chat with the travel agent
func Chat(req ChatRequest) (*AIChatResponse, error) {
	client := GetAIClient()

	// Convert request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chat request: %w", err)
	}

	// Make HTTP request to LangGraph agent
	resp, err := client.makeRequest("POST", "/chat", jsonData)
	if err != nil {
		return nil, err
	}

	// Parse response
	var chatResp AIChatResponse
	if err := json.Unmarshal(resp, &chatResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal chat response: %w", err)
	}

	return &chatResp, nil
}

// GenerateAIPackingList generates a personalized packing list using the LangGraph agent
func GenerateAIPackingList(req AIPackingRequest) (*AIPackingResponse, error) {
	client := GetAIClient()

	// Convert request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal packing request: %w", err)
	}

	// Make HTTP request to LangGraph agent
	resp, err := client.makeRequest("POST", "/generate-packing-list", jsonData)
	if err != nil {
		return nil, err
	}

	// Parse response
	var packingResp AIPackingResponse
	if err := json.Unmarshal(resp, &packingResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal packing response: %w", err)
	}

	return &packingResp, nil
}

// GetAIRecommendations gets recommendations for a specific city and category
func GetAIRecommendations(city, category string) ([]map[string]interface{}, error) {
	client := GetAIClient()

	// Build URL with query parameters
	url := fmt.Sprintf("%s/tools/recommendations?city=%s&category=%s", client.baseURL, city, category)

	// Make HTTP request
	resp, err := client.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get recommendations: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("recommendations API returned status: %d", resp.StatusCode)
	}

	// Parse response
	var result struct {
		Success         bool                     `json:"success"`
		Recommendations []map[string]interface{} `json:"recommendations"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode recommendations response: %w", err)
	}

	return result.Recommendations, nil
}

// GetAIEvents gets events for a specific city and date
func GetAIEvents(city, date string) ([]map[string]interface{}, error) {
	client := GetAIClient()

	// Build URL with query parameters
	url := fmt.Sprintf("%s/tools/events?city=%s", client.baseURL, city)
	if date != "" {
		url += fmt.Sprintf("&date=%s", date)
	}

	// Make HTTP request
	resp, err := client.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("events API returned status: %d", resp.StatusCode)
	}

	// Parse response
	var result struct {
		Success bool                     `json:"success"`
		Events  []map[string]interface{} `json:"events"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode events response: %w", err)
	}

	return result.Events, nil
}

// GetAttractions gets attractions for a specific city and category
func GetAttractions(city, category string) ([]map[string]interface{}, error) {
	client := GetAIClient()

	// Build URL with query parameters
	url := fmt.Sprintf("%s/tools/attractions?city=%s&category=%s", client.baseURL, city, category)

	// Make HTTP request
	resp, err := client.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get attractions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("attractions API returned status: %d", resp.StatusCode)
	}

	// Parse response
	var result struct {
		Success     bool                     `json:"success"`
		Attractions []map[string]interface{} `json:"attractions"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode attractions response: %w", err)
	}

	return result.Attractions, nil
}

// HealthCheck checks if the LangGraph agent is healthy
func HealthCheck() error {
	client := GetAIClient()

	resp, err := client.httpClient.Get(client.baseURL + "/health")
	if err != nil {
		return fmt.Errorf("failed to check health: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
	}

	return nil
}

// makeRequest makes an HTTP request to the LangGraph agent
func (c *AIClient) makeRequest(method, endpoint string, data []byte) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	// Read response
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return responseData, nil
}

// Legacy functions for backward compatibility
func GenerateItineraryLegacy(req interface{}) (interface{}, error) {
	// Convert legacy request to new format
	itineraryReq, ok := req.(ItineraryRequest)
	if !ok {
		return nil, fmt.Errorf("invalid request type")
	}

	result, err := GenerateItinerary(itineraryReq)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// SaveItinerary saves an itinerary (placeholder for future implementation)
func SaveItinerary(itinerary interface{}) error {
	// TODO: Implement saving to GCS or database
	return nil
}

// GetItinerary retrieves an itinerary by ID (placeholder for future implementation)
func GetItinerary(id string) (interface{}, error) {
	// TODO: Implement retrieval from GCS or database
	return map[string]interface{}{
		"id":         id,
		"city":       "Sample City",
		"total_cost": 1000.0,
		"duration":   7,
		"days":       []interface{}{},
		"created_at": "2024-01-01T00:00:00Z",
	}, nil
}

// DeleteItinerary deletes an itinerary (placeholder for future implementation)
func DeleteItinerary(id string) error {
	// TODO: Implement deletion from GCS or database
	return nil
}
