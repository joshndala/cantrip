package services

//COMPLETED
import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// Event represents an event in a city
type Event struct {
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Date             string   `json:"date"`
	EndDate          string   `json:"end_date,omitempty"`
	Time             string   `json:"time,omitempty"`
	Location         string   `json:"location"`
	Price            float64  `json:"price"`
	PriceRange       string   `json:"price_range,omitempty"`
	Category         string   `json:"category"`
	Type             string   `json:"type,omitempty"`
	TicketsAvailable bool     `json:"tickets_available"`
	BookingURL       string   `json:"booking_url,omitempty"`
	Rating           float64  `json:"rating,omitempty"`
	Tags             []string `json:"tags,omitempty"`
}

// TripSuggestion represents a trip suggestion
type TripSuggestion struct {
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Activities    []string `json:"activities"`
	EstimatedCost float64  `json:"estimated_cost"`
	Duration      int      `json:"duration"`
	Tags          []string `json:"tags"`
}

// EventAPIResponse represents the response from event APIs
type EventAPIResponse struct {
	Events []Event `json:"events"`
}

// MoodInterests maps moods to relevant event categories and interests
var MoodInterests = map[string][]string{
	"excited":     {"music", "sports", "festival", "entertainment"},
	"relaxed":     {"arts", "culture", "museum", "theater"},
	"adventurous": {"outdoor", "sports", "adventure", "festival"},
	"romantic":    {"arts", "music", "dining", "theater"},
	"family":      {"family", "kids", "entertainment", "outdoor"},
	"cultural":    {"culture", "arts", "museum", "heritage"},
	"party":       {"music", "nightlife", "festival", "entertainment"},
	"educational": {"museum", "arts", "culture", "workshop"},
}

// GetEvents retrieves events for a city based on mood and interests
func GetEvents(city, mood string, interests []string) ([]Event, error) {
	// First, try to get events from real APIs
	if events, err := getEventsFromAPI(city, mood, interests); err == nil && len(events) > 0 {
		return events, nil
	}

	// Fallback to sample event data
	return getEventsFromSampleData(city, mood, interests)
}

// FilterEventsByDate filters events by a specific date
func FilterEventsByDate(events []Event, targetDate string) []Event {
	var filteredEvents []Event
	
	for _, event := range events {
		// Check if event is on the target date or within a range
		if event.Date == targetDate {
			filteredEvents = append(filteredEvents, event)
		} else if event.EndDate != "" {
			// Check if target date falls within the event's date range
			if targetDate >= event.Date && targetDate <= event.EndDate {
				filteredEvents = append(filteredEvents, event)
			}
		}
	}
	
	return filteredEvents
}

// getEventsFromAPI attempts to get events from real event APIs
func getEventsFromAPI(city, mood string, interests []string) ([]Event, error) {
	// Check if API keys are configured
	ticketmasterKey := os.Getenv("TICKETMASTER_API_KEY")
	eventbriteKey := os.Getenv("EVENTBRITE_API_KEY")

	if ticketmasterKey == "" && eventbriteKey == "" {
		return nil, fmt.Errorf("no event API keys configured")
	}

	var allEvents []Event

	// Try Ticketmaster API
	if ticketmasterKey != "" {
		if events, err := getTicketmasterEvents(city, mood, interests); err == nil {
			allEvents = append(allEvents, events...)
		}
	}

	// Try Eventbrite API
	if eventbriteKey != "" {
		if events, err := getEventbriteEvents(city, mood, interests); err == nil {
			allEvents = append(allEvents, events...)
		}
	}

	// Filter and rank events based on mood and interests
	filteredEvents := filterEventsByMoodAndInterests(allEvents, mood, interests)

	return filteredEvents, nil
}

// getTicketmasterEvents gets events from Ticketmaster API
func getTicketmasterEvents(city, mood string, interests []string) ([]Event, error) {
	apiKey := os.Getenv("TICKETMASTER_API_KEY")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Build query parameters
	params := fmt.Sprintf("?apikey=%s&city=%s&size=20", apiKey, city)

	url := fmt.Sprintf("https://app.ticketmaster.com/discovery/v2/events.json%s", params)

	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Ticketmaster request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Ticketmaster events: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ticketmaster API returned status: %d", resp.StatusCode)
	}

	// Parse Ticketmaster response (simplified)
	var apiResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode Ticketmaster response: %w", err)
	}

	// Convert Ticketmaster response to our Event format
	events := convertTicketmasterResponse(apiResponse)

	return events, nil
}

// getEventbriteEvents gets events from Eventbrite API
func getEventbriteEvents(city, mood string, interests []string) ([]Event, error) {
	apiKey := os.Getenv("EVENTBRITE_API_KEY")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Build query parameters
	params := fmt.Sprintf("?location.address=%s&expand=venue", city)

	url := fmt.Sprintf("https://www.eventbriteapi.com/v3/events/search/%s", params)

	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Eventbrite request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Eventbrite events: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Eventbrite API returned status: %d", resp.StatusCode)
	}

	// Parse Eventbrite response (simplified)
	var apiResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode Eventbrite response: %w", err)
	}

	// Convert Eventbrite response to our Event format
	events := convertEventbriteResponse(apiResponse)

	return events, nil
}

// convertTicketmasterResponse converts Ticketmaster API response to our Event format
func convertTicketmasterResponse(response map[string]interface{}) []Event {
	var events []Event

	if embedded, ok := response["_embedded"].(map[string]interface{}); ok {
		if eventsList, ok := embedded["events"].([]interface{}); ok {
			for _, eventData := range eventsList {
				if eventMap, ok := eventData.(map[string]interface{}); ok {
					event := Event{
						Name:             getString(eventMap, "name"),
						Description:      getString(eventMap, "description"),
						Date:             getString(eventMap, "dates.start.localDate"),
						Time:             getString(eventMap, "dates.start.localTime"),
						Location:         getString(eventMap, "_embedded.venues.0.name"),
						Price:            getPriceFromTicketmaster(eventMap),
						Category:         getString(eventMap, "classifications.0.segment.name"),
						Type:             getString(eventMap, "classifications.0.genre.name"),
						TicketsAvailable: true,
						BookingURL:       getString(eventMap, "url"),
						Rating:           4.0, // Default rating
						Tags:             getTagsFromTicketmaster(eventMap),
					}
					events = append(events, event)
				}
			}
		}
	}

	return events
}

// convertEventbriteResponse converts Eventbrite API response to our Event format
func convertEventbriteResponse(response map[string]interface{}) []Event {
	var events []Event

	if eventsList, ok := response["events"].([]interface{}); ok {
		for _, eventData := range eventsList {
			if eventMap, ok := eventData.(map[string]interface{}); ok {
				event := Event{
					Name:             getString(eventMap, "name.text"),
					Description:      getString(eventMap, "description.text"),
					Date:             getString(eventMap, "start.local"),
					EndDate:          getString(eventMap, "end.local"),
					Location:         getString(eventMap, "venue.name"),
					Price:            getPriceFromEventbrite(eventMap),
					Category:         getString(eventMap, "category.name"),
					TicketsAvailable: true,
					BookingURL:       getString(eventMap, "url"),
					Rating:           4.0, // Default rating
					Tags:             getTagsFromEventbrite(eventMap),
				}
				events = append(events, event)
			}
		}
	}

	return events
}

// getEventsFromSampleData gets events from sample data based on mood and interests
func getEventsFromSampleData(city, mood string, interests []string) ([]Event, error) {
	// Load city metadata to get real attractions and activities
	metadata, err := loadCityMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to load city metadata: %w", err)
	}

	// Find the city in metadata
	cityData, err := findCity(metadata, city)
	if err != nil {
		// Return generic events if city not found
		return getGenericEvents(city, mood, interests), nil
	}

	// Generate events based on city's attractions and seasonal activities
	events := generateEventsFromCityData(cityData, mood, interests)

	// Filter events based on mood and interests
	filteredEvents := filterEventsByMoodAndInterests(events, mood, interests)

	return filteredEvents, nil
}

// generateEventsFromCityData creates events based on city metadata
func generateEventsFromCityData(cityData *City, mood string, interests []string) []Event {
	var events []Event

	// Get current season for relevant activities
	currentSeason := getCurrentSeason()
	seasonData, exists := cityData.Seasons[currentSeason]

	// Create events from attractions
	for _, attraction := range cityData.Attractions {
		event := Event{
			Name:             fmt.Sprintf("Visit %s", attraction),
			Description:      fmt.Sprintf("Explore %s in %s", attraction, cityData.Name),
			Date:             "", // No specific date - user can check availability
			Location:         fmt.Sprintf("%s, %s", attraction, cityData.Name),
			Price:            0, // No pricing info available
			Category:         "attraction",
			Type:             "sightseeing",
			TicketsAvailable: false, // Unknown availability
			Rating:           4.0,   // Default rating
			Tags:             []string{"attraction", "sightseeing", "tourism"},
		}
		events = append(events, event)
	}

	// Create events from seasonal activities
	if exists {
		for _, activity := range seasonData.Activities {
			event := Event{
				Name:             activity,
				Description:      fmt.Sprintf("Experience %s in %s during %s", activity, cityData.Name, currentSeason),
				Date:             "", // Seasonal activity - check local schedules
				Location:         cityData.Name,
				Price:            0, // No pricing info available
				Category:         "activity",
				Type:             "seasonal",
				TicketsAvailable: false, // Unknown availability
				Rating:           4.0,   // Default rating
				Tags:             []string{"activity", currentSeason, "local"},
			}
			events = append(events, event)
		}
	}

	// Add some neighborhood exploration events
	for _, neighborhood := range cityData.Neighborhoods {
		event := Event{
			Name:             fmt.Sprintf("Explore %s", neighborhood),
			Description:      fmt.Sprintf("Discover the vibrant %s neighborhood in %s", neighborhood, cityData.Name),
			Date:             "", // Always available
			Location:         fmt.Sprintf("%s, %s", neighborhood, cityData.Name),
			Price:            0, // Free exploration
			Category:         "neighborhood",
			Type:             "exploration",
			TicketsAvailable: true, // Always available
			Rating:           4.0,  // Default rating
			Tags:             []string{"neighborhood", "local", "exploration"},
		}
		events = append(events, event)
	}

	return events
}

// getGenericEvents returns generic events for cities not in metadata
func getGenericEvents(city, mood string, interests []string) []Event {
	// Create generic events based on mood and interests
	var events []Event

	// Add some generic exploration events
	events = append(events, Event{
		Name:             fmt.Sprintf("Explore Downtown %s", city),
		Description:      fmt.Sprintf("Discover the heart of %s with its shops, restaurants, and attractions", city),
		Date:             "", // Always available
		Location:         fmt.Sprintf("Downtown %s", city),
		Price:            0, // Free exploration
		Category:         "exploration",
		Type:             "sightseeing",
		TicketsAvailable: true, // Always available
		Rating:           4.0,  // Default rating
		Tags:             []string{"downtown", "exploration", "local"},
	})

	// Add mood-based generic events
	moodCategories := MoodInterests[strings.ToLower(mood)]
	if moodCategories != nil {
		for _, category := range moodCategories {
			events = append(events, Event{
				Name:             fmt.Sprintf("Local %s Experience in %s", strings.Title(category), city),
				Description:      fmt.Sprintf("Experience the local %s scene in %s", category, city),
				Date:             "", // Check local schedules
				Location:         city,
				Price:            0, // No pricing info available
				Category:         category,
				Type:             "local",
				TicketsAvailable: false, // Unknown availability
				Rating:           4.0,   // Default rating
				Tags:             []string{category, "local", "experience"},
			})
		}
	}

	return events
}

// filterEventsByMoodAndInterests filters events based on mood and user interests
func filterEventsByMoodAndInterests(events []Event, mood string, interests []string) []Event {
	var filteredEvents []Event

	// Get mood-based categories
	moodCategories := MoodInterests[strings.ToLower(mood)]
	if moodCategories == nil {
		moodCategories = []string{"entertainment"} // Default category
	}

	// Create a combined list of interests
	allInterests := append(interests, moodCategories...)

	for _, event := range events {
		// Check if event matches any interest or mood category
		if matchesInterests(event, allInterests) {
			filteredEvents = append(filteredEvents, event)
		}
	}

	// Sort by rating (highest first)
	sortEventsByRating(filteredEvents)

	// Limit to top 10 events
	if len(filteredEvents) > 10 {
		filteredEvents = filteredEvents[:10]
	}

	return filteredEvents
}

// matchesInterests checks if an event matches any of the given interests
func matchesInterests(event Event, interests []string) bool {
	eventText := strings.ToLower(event.Name + " " + event.Description + " " + event.Category + " " + event.Type)

	for _, tag := range event.Tags {
		eventText += " " + strings.ToLower(tag)
	}

	for _, interest := range interests {
		if strings.Contains(eventText, strings.ToLower(interest)) {
			return true
		}
	}

	return false
}

// sortEventsByRating sorts events by rating (highest first)
func sortEventsByRating(events []Event) {
	// Simple bubble sort for small lists
	for i := 0; i < len(events)-1; i++ {
		for j := 0; j < len(events)-i-1; j++ {
			if events[j].Rating < events[j+1].Rating {
				events[j], events[j+1] = events[j+1], events[j]
			}
		}
	}
}

// Helper functions for API response parsing
func getString(data map[string]interface{}, path string) string {
	keys := strings.Split(path, ".")
	current := data

	for _, key := range keys {
		if val, ok := current[key].(map[string]interface{}); ok {
			current = val
		} else if val, ok := current[key].(string); ok {
			return val
		} else {
			return ""
		}
	}

	return ""
}

func getPriceFromTicketmaster(data map[string]interface{}) float64 {
	// Simplified price extraction from Ticketmaster response
	if priceRanges, ok := data["priceRanges"].([]interface{}); ok && len(priceRanges) > 0 {
		if priceRange, ok := priceRanges[0].(map[string]interface{}); ok {
			if min, ok := priceRange["min"].(float64); ok {
				return min
			}
		}
	}
	return 25.0 // Default price
}

func getPriceFromEventbrite(data map[string]interface{}) float64 {
	// Simplified price extraction from Eventbrite response
	if ticketClasses, ok := data["ticket_availability"].(map[string]interface{}); ok {
		if price, ok := ticketClasses["minimum_ticket_price"].(map[string]interface{}); ok {
			if value, ok := price["value"].(float64); ok {
				return value / 100 // Convert cents to dollars
			}
		}
	}
	return 25.0 // Default price
}

func getTagsFromTicketmaster(data map[string]interface{}) []string {
	var tags []string
	if classifications, ok := data["classifications"].([]interface{}); ok && len(classifications) > 0 {
		if classification, ok := classifications[0].(map[string]interface{}); ok {
			if segment, ok := classification["segment"].(map[string]interface{}); ok {
				if name, ok := segment["name"].(string); ok {
					tags = append(tags, strings.ToLower(name))
				}
			}
		}
	}
	return tags
}

func getTagsFromEventbrite(data map[string]interface{}) []string {
	var tags []string
	if category, ok := data["category"].(map[string]interface{}); ok {
		if name, ok := category["name"].(string); ok {
			tags = append(tags, strings.ToLower(name))
		}
	}
	return tags
}

// GenerateTripSuggestions generates trip suggestions based on mood and interests
func GenerateTripSuggestions(mood, city string, budget float64, duration int, interests []string, weather WeatherInfo) ([]TripSuggestion, error) {
	// Load city metadata to get real attractions and activities
	metadata, err := loadCityMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to load city metadata: %w", err)
	}

	// Find the city in metadata
	cityData, err := findCity(metadata, city)
	if err != nil {
		// Generate generic suggestions if city not found
		return generateGenericTripSuggestions(mood, city, budget, duration, interests, weather), nil
	}

	// Generate suggestions based on city data
	suggestions := generateCityBasedTripSuggestions(cityData, mood, budget, duration, interests, weather)

	return suggestions, nil
}

// generateCityBasedTripSuggestions creates trip suggestions based on city metadata
func generateCityBasedTripSuggestions(cityData *City, mood string, budget float64, duration int, interests []string, weather WeatherInfo) []TripSuggestion {
	var suggestions []TripSuggestion

	// Get current season for relevant activities
	currentSeason := getCurrentSeason()
	seasonData, exists := cityData.Seasons[currentSeason]

	// 1. Cultural Explorer Suggestion
	suggestions = append(suggestions, TripSuggestion{
		Title:         fmt.Sprintf("Cultural Explorer in %s", cityData.Name),
		Description:   fmt.Sprintf("Immerse yourself in the rich culture of %s with museums, galleries, and historic sites", cityData.Name),
		Activities:    getCulturalActivities(cityData, &seasonData),
		EstimatedCost: calculateCulturalCost(budget, duration),
		Duration:      duration,
		Tags:          []string{"culture", "arts", "history", "museum"},
	})

	// 2. Outdoor Adventure Suggestion
	if isGoodWeatherForOutdoor(weather) {
		suggestions = append(suggestions, TripSuggestion{
			Title:         fmt.Sprintf("Outdoor Adventure in %s", cityData.Name),
			Description:   fmt.Sprintf("Explore the natural beauty and outdoor activities in %s", cityData.Name),
			Activities:    getOutdoorActivities(cityData, &seasonData),
			EstimatedCost: calculateOutdoorCost(budget, duration),
			Duration:      duration,
			Tags:          []string{"outdoor", "nature", "adventure", "active"},
		})
	}

	// 3. Food & Local Experience Suggestion
	suggestions = append(suggestions, TripSuggestion{
		Title:         fmt.Sprintf("Local Food & Culture in %s", cityData.Name),
		Description:   fmt.Sprintf("Taste the local cuisine and experience the authentic %s lifestyle", cityData.Name),
		Activities:    getFoodAndLocalActivities(cityData, &seasonData),
		EstimatedCost: calculateFoodCost(budget, duration),
		Duration:      duration,
		Tags:          []string{"food", "local", "culture", "dining"},
	})

	// 4. Neighborhood Explorer Suggestion
	suggestions = append(suggestions, TripSuggestion{
		Title:         fmt.Sprintf("Neighborhood Explorer in %s", cityData.Name),
		Description:   fmt.Sprintf("Discover the diverse neighborhoods and local life in %s", cityData.Name),
		Activities:    getNeighborhoodActivities(cityData),
		EstimatedCost: calculateNeighborhoodCost(budget, duration),
		Duration:      duration,
		Tags:          []string{"neighborhood", "local", "exploration", "community"},
	})

	// 5. Seasonal Special Suggestion
	if exists {
		suggestions = append(suggestions, TripSuggestion{
			Title:         fmt.Sprintf("%s Seasonal Experience in %s", strings.Title(currentSeason), cityData.Name),
			Description:   fmt.Sprintf("Experience the best of %s during %s with seasonal activities and events", cityData.Name, currentSeason),
			Activities:    seasonData.Activities,
			EstimatedCost: calculateSeasonalCost(budget, duration),
			Duration:      duration,
			Tags:          append([]string{currentSeason, "seasonal"}, interests...),
		})
	}

	// 6. Budget-Friendly Suggestion
	suggestions = append(suggestions, TripSuggestion{
		Title:         fmt.Sprintf("Budget-Friendly %s Experience", cityData.Name),
		Description:   fmt.Sprintf("Explore %s on a budget with free and low-cost activities", cityData.Name),
		Activities:    getBudgetActivities(cityData, &seasonData),
		EstimatedCost: budget * 0.6, // 60% of budget
		Duration:      duration,
		Tags:          []string{"budget", "affordable", "free", "value"},
	})

	// Filter suggestions based on mood and interests
	filteredSuggestions := filterSuggestionsByMoodAndInterests(suggestions, mood, interests)

	// Limit to top 5 suggestions
	if len(filteredSuggestions) > 5 {
		filteredSuggestions = filteredSuggestions[:5]
	}

	return filteredSuggestions
}

// generateGenericTripSuggestions creates generic suggestions for cities not in metadata
func generateGenericTripSuggestions(mood, city string, budget float64, duration int, interests []string, weather WeatherInfo) []TripSuggestion {
	var suggestions []TripSuggestion

	// Generic cultural suggestion
	suggestions = append(suggestions, TripSuggestion{
		Title:         fmt.Sprintf("Discover %s", city),
		Description:   fmt.Sprintf("Explore the culture, history, and attractions of %s", city),
		Activities:    []string{"Visit local museums", "Explore downtown", "Try local cuisine", "Visit historic sites"},
		EstimatedCost: budget * 0.8,
		Duration:      duration,
		Tags:          append([]string{"culture", "exploration"}, interests...),
	})

	// Mood-based suggestion
	moodCategories := MoodInterests[strings.ToLower(mood)]
	if moodCategories != nil {
		activities := []string{}
		for _, category := range moodCategories {
			activities = append(activities, fmt.Sprintf("Experience local %s", category))
		}
		activities = append(activities, "Explore the city center", "Try local restaurants")

		suggestions = append(suggestions, TripSuggestion{
			Title:         fmt.Sprintf("%s Adventure in %s", strings.Title(mood), city),
			Description:   fmt.Sprintf("Enjoy a %s experience in %s with activities tailored to your mood", mood, city),
			Activities:    activities,
			EstimatedCost: budget * 0.9,
			Duration:      duration,
			Tags:          append(moodCategories, interests...),
		})
	}

	return suggestions
}

// Helper functions for generating activities
func getCulturalActivities(cityData *City, seasonData *Season) []string {
	var activities []string

	// Add museums and cultural attractions
	for _, attraction := range cityData.Attractions {
		if strings.Contains(strings.ToLower(attraction), "museum") ||
			strings.Contains(strings.ToLower(attraction), "gallery") ||
			strings.Contains(strings.ToLower(attraction), "art") {
			activities = append(activities, fmt.Sprintf("Visit %s", attraction))
		}
	}

	// Add cultural activities from seasonal data
	if seasonData != nil {
		for _, activity := range seasonData.Activities {
			if strings.Contains(strings.ToLower(activity), "festival") ||
				strings.Contains(strings.ToLower(activity), "cultural") ||
				strings.Contains(strings.ToLower(activity), "arts") {
				activities = append(activities, activity)
			}
		}
	}

	// Add generic cultural activities
	activities = append(activities, "Explore historic districts", "Visit local galleries", "Attend cultural events")

	return activities
}

func getOutdoorActivities(cityData *City, seasonData *Season) []string {
	var activities []string

	// Add outdoor attractions
	for _, attraction := range cityData.Attractions {
		if strings.Contains(strings.ToLower(attraction), "park") ||
			strings.Contains(strings.ToLower(attraction), "mountain") ||
			strings.Contains(strings.ToLower(attraction), "beach") ||
			strings.Contains(strings.ToLower(attraction), "island") {
			activities = append(activities, fmt.Sprintf("Explore %s", attraction))
		}
	}

	// Add seasonal outdoor activities
	if seasonData != nil {
		for _, activity := range seasonData.Activities {
			if strings.Contains(strings.ToLower(activity), "hiking") ||
				strings.Contains(strings.ToLower(activity), "beach") ||
				strings.Contains(strings.ToLower(activity), "skiing") ||
				strings.Contains(strings.ToLower(activity), "outdoor") {
				activities = append(activities, activity)
			}
		}
	}

	// Add generic outdoor activities
	activities = append(activities, "Take walking tours", "Visit parks and gardens", "Enjoy outdoor dining")

	return activities
}

func getFoodAndLocalActivities(cityData *City, seasonData *Season) []string {
	var activities []string

	// Add food-related seasonal activities
	if seasonData != nil {
		for _, activity := range seasonData.Activities {
			if strings.Contains(strings.ToLower(activity), "food") ||
				strings.Contains(strings.ToLower(activity), "culinary") ||
				strings.Contains(strings.ToLower(activity), "wine") {
				activities = append(activities, activity)
			}
		}
	}

	// Add food markets and local experiences
	for _, attraction := range cityData.Attractions {
		if strings.Contains(strings.ToLower(attraction), "market") ||
			strings.Contains(strings.ToLower(attraction), "district") {
			activities = append(activities, fmt.Sprintf("Visit %s", attraction))
		}
	}

	// Add generic food activities
	activities = append(activities, "Try local restaurants", "Visit food markets", "Take cooking classes", "Sample local cuisine")

	return activities
}

func getNeighborhoodActivities(cityData *City) []string {
	var activities []string

	// Add neighborhood exploration
	for _, neighborhood := range cityData.Neighborhoods {
		activities = append(activities, fmt.Sprintf("Explore %s", neighborhood))
	}

	// Add generic neighborhood activities
	activities = append(activities, "Walk through local markets", "Visit neighborhood cafes", "Experience local nightlife", "Shop at local boutiques")

	return activities
}

func getBudgetActivities(cityData *City, seasonData *Season) []string {
	var activities []string

	// Add free attractions
	for _, attraction := range cityData.Attractions {
		if strings.Contains(strings.ToLower(attraction), "park") ||
			strings.Contains(strings.ToLower(attraction), "market") ||
			strings.Contains(strings.ToLower(attraction), "district") {
			activities = append(activities, fmt.Sprintf("Visit %s", attraction))
		}
	}

	// Add free seasonal activities
	if seasonData != nil {
		for _, activity := range seasonData.Activities {
			if strings.Contains(strings.ToLower(activity), "walk") ||
				strings.Contains(strings.ToLower(activity), "hiking") ||
				strings.Contains(strings.ToLower(activity), "viewing") {
				activities = append(activities, activity)
			}
		}
	}

	// Add generic budget activities
	activities = append(activities, "Take free walking tours", "Visit public parks", "Explore on foot", "Window shop in local areas")

	return activities
}

// Helper functions for cost calculation
func calculateCulturalCost(budget float64, duration int) float64 {
	return budget * 0.7 // Cultural activities are typically moderate cost
}

func calculateOutdoorCost(budget float64, duration int) float64 {
	return budget * 0.5 // Outdoor activities are often lower cost
}

func calculateFoodCost(budget float64, duration int) float64 {
	return budget * 0.8 // Food experiences can be expensive
}

func calculateNeighborhoodCost(budget float64, duration int) float64 {
	return budget * 0.6 // Neighborhood exploration is moderate cost
}

func calculateSeasonalCost(budget float64, duration int) float64 {
	return budget * 0.75 // Seasonal activities vary in cost
}

// Helper functions for filtering
func isGoodWeatherForOutdoor(weather WeatherInfo) bool {
	// Consider weather suitable for outdoor activities
	return weather.Temperature > 10 &&
		weather.Temperature < 35 &&
		!strings.Contains(strings.ToLower(weather.Condition), "rain") &&
		!strings.Contains(strings.ToLower(weather.Condition), "snow")
}

func filterSuggestionsByMoodAndInterests(suggestions []TripSuggestion, mood string, interests []string) []TripSuggestion {
	var filteredSuggestions []TripSuggestion

	// Get mood-based categories
	moodCategories := MoodInterests[strings.ToLower(mood)]
	if moodCategories == nil {
		moodCategories = []string{"entertainment"} // Default category
	}

	// Create a combined list of interests
	allInterests := append(interests, moodCategories...)

	for _, suggestion := range suggestions {
		// Check if suggestion matches any interest or mood category
		if matchesSuggestionInterests(suggestion, allInterests) {
			filteredSuggestions = append(filteredSuggestions, suggestion)
		}
	}

	return filteredSuggestions
}

func matchesSuggestionInterests(suggestion TripSuggestion, interests []string) bool {
	suggestionText := strings.ToLower(suggestion.Title + " " + suggestion.Description)

	for _, tag := range suggestion.Tags {
		suggestionText += " " + strings.ToLower(tag)
	}

	for _, interest := range interests {
		if strings.Contains(suggestionText, strings.ToLower(interest)) {
			return true
		}
	}

	return false
}
