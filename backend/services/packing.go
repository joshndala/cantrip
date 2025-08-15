package services

//COMPLETED
import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type PackingRequest struct {
	Destination  string   `json:"destination"`
	StartDate    string   `json:"start_date"`
	EndDate      string   `json:"end_date"`
	Activities   []string `json:"activities"`
	Weather      string   `json:"weather"`
	GroupSize    int      `json:"group_size"`
	AgeGroup     string   `json:"age_group"`
	SpecialNeeds []string `json:"special_needs"`
	BaggageType  string   `json:"baggage_type"`
}

type PackingResponse struct {
	ID          string        `json:"id"`
	Destination string        `json:"destination"`
	Categories  []interface{} `json:"categories"`
	TotalItems  int           `json:"total_items"`
	Notes       []string      `json:"notes"`
	Weather     WeatherInfo   `json:"weather"`
}

// PackingCategory represents a category of items in the packing list
type PackingCategory struct {
	Name  string        `json:"name"`
	Items []PackingItem `json:"items"`
}

// PackingItem represents a single item in the packing list
type PackingItem struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Reason   string `json:"reason"`
}

// PackingRules represents the structure of packing_rules.json
type PackingRules struct {
	WeatherRules  map[string]interface{} `json:"weather_rules"`
	ActivityRules map[string]interface{} `json:"activity_rules"`
	DurationRules map[string]interface{} `json:"duration_rules"`
	GroupRules    map[string]interface{} `json:"group_rules"`
	AgeRules      map[string]interface{} `json:"age_rules"`
	SpecialNeeds  map[string]interface{} `json:"special_needs"`
	BaggageRules  map[string]interface{} `json:"baggage_rules"`
}

// GeneratePackingList generates a packing list based on the request and weather information
func GeneratePackingList(req PackingRequest, weather WeatherInfo) (PackingResponse, error) {
	// Load packing rules
	rules, err := loadPackingRules()
	if err != nil {
		return PackingResponse{}, fmt.Errorf("failed to load packing rules: %w", err)
	}

	// Calculate trip duration
	duration, err := calculateDuration(req.StartDate, req.EndDate)
	if err != nil {
		return PackingResponse{}, fmt.Errorf("failed to calculate duration: %w", err)
	}

	// Determine weather category based on temperature
	weatherCategory := getWeatherCategory(weather.Temperature)

	// Generate categories based on weather, activities, and other factors
	categories := []PackingCategory{}

	// Add weather-based clothing
	if weatherItems := getWeatherItems(rules, weatherCategory); len(weatherItems) > 0 {
		categories = append(categories, PackingCategory{
			Name:  "Weather-Appropriate Clothing",
			Items: weatherItems,
		})
	}

	// Add activity-based items
	for _, activity := range req.Activities {
		if activityItems := getActivityItems(rules, activity); len(activityItems) > 0 {
			categories = append(categories, PackingCategory{
				Name:  fmt.Sprintf("%s Gear", strings.Title(strings.ReplaceAll(activity, "_", " "))),
				Items: activityItems,
			})
		}
	}

	// Add age-specific items
	if ageItems := getAgeItems(rules, req.AgeGroup); len(ageItems) > 0 {
		categories = append(categories, PackingCategory{
			Name:  "Age-Specific Items",
			Items: ageItems,
		})
	}

	// Add special needs items
	for _, need := range req.SpecialNeeds {
		if specialItems := getSpecialNeedsItems(rules, need); len(specialItems) > 0 {
			categories = append(categories, PackingCategory{
				Name:  fmt.Sprintf("%s Items", strings.Title(strings.ReplaceAll(need, "_", " "))),
				Items: specialItems,
			})
		}
	}

	// Add essentials
	essentials := getEssentials(rules, req.BaggageType)
	if len(essentials) > 0 {
		categories = append(categories, PackingCategory{
			Name:  "Essentials",
			Items: essentials,
		})
	}

	// Apply duration multiplier
	applyDurationMultiplier(categories, rules, duration)

	// Apply group size multiplier
	applyGroupMultiplier(categories, rules, req.GroupSize)

	// Calculate total items
	totalItems := 0
	for _, category := range categories {
		for _, item := range category.Items {
			totalItems += item.Quantity
		}
	}

	// Generate notes
	notes := generateNotes(rules, duration, req.GroupSize, weatherCategory)

	// Convert categories to interface{} for JSON serialization
	categoriesInterface := make([]interface{}, len(categories))
	for i, category := range categories {
		categoriesInterface[i] = category
	}

	return PackingResponse{
		ID:          generatePackingListID(req.Destination, req.StartDate),
		Destination: req.Destination,
		Categories:  categoriesInterface,
		TotalItems:  totalItems,
		Notes:       notes,
		Weather:     weather,
	}, nil
}

// loadPackingRules loads the packing rules from the JSON file
func loadPackingRules() (*PackingRules, error) {
	data, err := os.ReadFile("data/packing_rules.json")
	if err != nil {
		return nil, err
	}

	var rules PackingRules
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, err
	}

	return &rules, nil
}

// calculateDuration calculates the duration of the trip in days
func calculateDuration(startDate, endDate string) (int, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return 0, err
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return 0, err
	}

	duration := end.Sub(start)
	return int(duration.Hours() / 24), nil
}

// getWeatherCategory determines the weather category based on temperature
func getWeatherCategory(temperature float64) string {
	switch {
	case temperature >= 25:
		return "hot"
	case temperature >= 15:
		return "warm"
	case temperature >= 5:
		return "mild"
	case temperature >= -5:
		return "cool"
	default:
		return "cold"
	}
}

// getWeatherItems gets items based on weather category
func getWeatherItems(rules *PackingRules, weatherCategory string) []PackingItem {
	var items []PackingItem

	if weatherRule, exists := rules.WeatherRules[weatherCategory]; exists {
		if weatherMap, ok := weatherRule.(map[string]interface{}); ok {
			// Add clothing
			if clothing, ok := weatherMap["clothing"].([]interface{}); ok {
				for _, item := range clothing {
					items = append(items, PackingItem{
						Name:     item.(string),
						Quantity: 1,
						Reason:   fmt.Sprintf("Appropriate for %s weather", weatherCategory),
					})
				}
			}

			// Add accessories
			if accessories, ok := weatherMap["accessories"].([]interface{}); ok {
				for _, item := range accessories {
					items = append(items, PackingItem{
						Name:     item.(string),
						Quantity: 1,
						Reason:   fmt.Sprintf("Essential for %s weather", weatherCategory),
					})
				}
			}

			// Add footwear
			if footwear, ok := weatherMap["footwear"].([]interface{}); ok {
				for _, item := range footwear {
					items = append(items, PackingItem{
						Name:     item.(string),
						Quantity: 1,
						Reason:   fmt.Sprintf("Suitable for %s weather", weatherCategory),
					})
				}
			}
		}
	}

	return items
}

// getActivityItems gets items based on activities
func getActivityItems(rules *PackingRules, activity string) []PackingItem {
	var items []PackingItem

	if activityRule, exists := rules.ActivityRules[activity]; exists {
		if activityMap, ok := activityRule.(map[string]interface{}); ok {
			// Add clothing
			if clothing, ok := activityMap["clothing"].([]interface{}); ok {
				for _, item := range clothing {
					items = append(items, PackingItem{
						Name:     item.(string),
						Quantity: 1,
						Reason:   fmt.Sprintf("Required for %s", activity),
					})
				}
			}

			// Add accessories
			if accessories, ok := activityMap["accessories"].([]interface{}); ok {
				for _, item := range accessories {
					items = append(items, PackingItem{
						Name:     item.(string),
						Quantity: 1,
						Reason:   fmt.Sprintf("Essential for %s", activity),
					})
				}
			}

			// Add footwear
			if footwear, ok := activityMap["footwear"].([]interface{}); ok {
				for _, item := range footwear {
					items = append(items, PackingItem{
						Name:     item.(string),
						Quantity: 1,
						Reason:   fmt.Sprintf("Suitable for %s", activity),
					})
				}
			}
		}
	}

	return items
}

// getAgeItems gets items based on age group
func getAgeItems(rules *PackingRules, ageGroup string) []PackingItem {
	var items []PackingItem

	if ageRule, exists := rules.AgeRules[ageGroup]; exists {
		if ageMap, ok := ageRule.(map[string]interface{}); ok {
			if additionalItems, ok := ageMap["additional_items"].([]interface{}); ok {
				for _, item := range additionalItems {
					items = append(items, PackingItem{
						Name:     item.(string),
						Quantity: 1,
						Reason:   fmt.Sprintf("Required for %s", ageGroup),
					})
				}
			}
		}
	}

	return items
}

// getSpecialNeedsItems gets items based on special needs
func getSpecialNeedsItems(rules *PackingRules, specialNeed string) []PackingItem {
	var items []PackingItem

	if specialRule, exists := rules.SpecialNeeds[specialNeed]; exists {
		if specialMap, ok := specialRule.(map[string]interface{}); ok {
			if additionalItems, ok := specialMap["additional_items"].([]interface{}); ok {
				for _, item := range additionalItems {
					items = append(items, PackingItem{
						Name:     item.(string),
						Quantity: 1,
						Reason:   fmt.Sprintf("Required for %s", specialNeed),
					})
				}
			}
		}
	}

	return items
}

// getEssentials gets essential items based on baggage type
func getEssentials(rules *PackingRules, baggageType string) []PackingItem {
	var items []PackingItem

	if baggageRule, exists := rules.BaggageRules[baggageType]; exists {
		if baggageMap, ok := baggageRule.(map[string]interface{}); ok {
			if essentials, ok := baggageMap["essentials"].([]interface{}); ok {
				for _, item := range essentials {
					items = append(items, PackingItem{
						Name:     item.(string),
						Quantity: 1,
						Reason:   "Essential item",
					})
				}
			}
		}
	}

	return items
}

// applyDurationMultiplier applies duration-based multipliers to item quantities
func applyDurationMultiplier(categories []PackingCategory, rules *PackingRules, duration int) {
	var durationCategory string
	switch {
	case duration <= 2:
		durationCategory = "weekend"
	case duration <= 7:
		durationCategory = "week"
	case duration <= 14:
		durationCategory = "two_weeks"
	default:
		durationCategory = "month"
	}

	if durationRule, exists := rules.DurationRules[durationCategory]; exists {
		if durationMap, ok := durationRule.(map[string]interface{}); ok {
			if multiplier, ok := durationMap["multiplier"].(float64); ok {
				for i := range categories {
					for j := range categories[i].Items {
						categories[i].Items[j].Quantity = int(math.Ceil(float64(categories[i].Items[j].Quantity) * multiplier))
					}
				}
			}
		}
	}
}

// applyGroupMultiplier applies group size-based multipliers to item quantities
func applyGroupMultiplier(categories []PackingCategory, rules *PackingRules, groupSize int) {
	var groupCategory string
	switch {
	case groupSize == 1:
		groupCategory = "solo"
	case groupSize == 2:
		groupCategory = "couple"
	case groupSize <= 4:
		groupCategory = "family"
	default:
		groupCategory = "group"
	}

	if groupRule, exists := rules.GroupRules[groupCategory]; exists {
		if groupMap, ok := groupRule.(map[string]interface{}); ok {
			if multiplier, ok := groupMap["multiplier"].(float64); ok {
				for i := range categories {
					for j := range categories[i].Items {
						categories[i].Items[j].Quantity = int(math.Ceil(float64(categories[i].Items[j].Quantity) * multiplier))
					}
				}
			}
		}
	}
}

// generateNotes generates helpful notes for the packing list
func generateNotes(rules *PackingRules, duration int, groupSize int, weatherCategory string) []string {
	var notes []string

	// Add duration note
	var durationCategory string
	switch {
	case duration <= 2:
		durationCategory = "weekend"
	case duration <= 7:
		durationCategory = "week"
	case duration <= 14:
		durationCategory = "two_weeks"
	default:
		durationCategory = "month"
	}

	if durationRule, exists := rules.DurationRules[durationCategory]; exists {
		if durationMap, ok := durationRule.(map[string]interface{}); ok {
			if note, ok := durationMap["notes"].(string); ok {
				notes = append(notes, note)
			}
		}
	}

	// Add group size note
	var groupCategory string
	switch {
	case groupSize == 1:
		groupCategory = "solo"
	case groupSize == 2:
		groupCategory = "couple"
	case groupSize <= 4:
		groupCategory = "family"
	default:
		groupCategory = "group"
	}

	if groupRule, exists := rules.GroupRules[groupCategory]; exists {
		if groupMap, ok := groupRule.(map[string]interface{}); ok {
			if note, ok := groupMap["notes"].(string); ok {
				notes = append(notes, note)
			}
		}
	}

	// Add weather note
	notes = append(notes, fmt.Sprintf("Weather is expected to be %s, pack accordingly", weatherCategory))

	return notes
}

// generatePackingListID generates a unique ID for the packing list
func generatePackingListID(destination, startDate string) string {
	return fmt.Sprintf("packing-%s-%s", strings.ToLower(strings.ReplaceAll(destination, " ", "-")), startDate)
}

// SavePackingList saves a packing list to GCS or local storage
func SavePackingList(packingList PackingResponse) error {
	// Try to save to GCS first
	if gcsClient := GetGCSClient(); gcsClient != nil {
		ctx := context.Background()
		objectName := fmt.Sprintf("packing_lists/%s.json", packingList.ID)

		if err := gcsClient.UploadJSON(ctx, objectName, packingList); err == nil {
			return nil
		}
		// If GCS fails, fall back to local storage
	}

	// Fallback to local storage
	return savePackingListLocal(packingList)
}

// savePackingListLocal saves a packing list to local file system
func savePackingListLocal(packingList PackingResponse) error {
	// Create the data directory if it doesn't exist
	dataDir := "data/packing_lists"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Create filename based on the packing list ID
	filename := filepath.Join(dataDir, fmt.Sprintf("%s.json", packingList.ID))

	// Convert the packing list to JSON
	jsonData, err := json.MarshalIndent(packingList, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal packing list to JSON: %w", err)
	}

	// Write the JSON data to file
	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write packing list to file: %w", err)
	}

	return nil
}

// GetPackingList retrieves a packing list by ID from GCS or local storage
func GetPackingList(id string) (PackingResponse, error) {
	// Try to get from GCS first
	if gcsClient := GetGCSClient(); gcsClient != nil {
		ctx := context.Background()
		objectName := fmt.Sprintf("packing_lists/%s.json", id)

		var packingList PackingResponse
		if err := gcsClient.DownloadJSON(ctx, objectName, &packingList); err == nil {
			return packingList, nil
		}
		// If GCS fails, fall back to local storage
	}

	// Fallback to local storage
	return getPackingListLocal(id)
}

// getPackingListLocal retrieves a packing list by ID from local file system
func getPackingListLocal(id string) (PackingResponse, error) {
	// Create the filename based on the packing list ID
	filename := filepath.Join("data/packing_lists", fmt.Sprintf("%s.json", id))

	// Check if the file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return PackingResponse{}, fmt.Errorf("packing list with ID '%s' not found", id)
	}

	// Read the JSON file
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return PackingResponse{}, fmt.Errorf("failed to read packing list file: %w", err)
	}

	// Unmarshal the JSON data into a PackingResponse struct
	var packingList PackingResponse
	if err := json.Unmarshal(jsonData, &packingList); err != nil {
		return PackingResponse{}, fmt.Errorf("failed to unmarshal packing list JSON: %w", err)
	}

	return packingList, nil
}

// GetPackingSuggestions gets packing suggestions based on destination, season, and activities
func GetPackingSuggestions(destination, season string, activities []string) ([]interface{}, error) {
	// Load packing rules
	rules, err := loadPackingRules()
	if err != nil {
		return nil, fmt.Errorf("failed to load packing rules: %w", err)
	}

	var suggestions []interface{}

	// Add weather-based suggestions based on season
	weatherCategory := getSeasonWeatherCategory(season)
	if weatherItems := getWeatherItems(rules, weatherCategory); len(weatherItems) > 0 {
		// Take a sample of weather items for suggestions
		sampleSize := min(5, len(weatherItems))
		for i := 0; i < sampleSize; i++ {
			suggestions = append(suggestions, map[string]interface{}{
				"item":     weatherItems[i].Name,
				"category": "weather",
				"reason":   fmt.Sprintf("Appropriate for %s weather during %s", weatherCategory, season),
			})
		}
	}

	// Add activity-based suggestions
	for _, activity := range activities {
		if activityItems := getActivityItems(rules, activity); len(activityItems) > 0 {
			// Take a sample of activity items for suggestions
			sampleSize := min(3, len(activityItems))
			for i := 0; i < sampleSize; i++ {
				suggestions = append(suggestions, map[string]interface{}{
					"item":     activityItems[i].Name,
					"category": "activity",
					"reason":   fmt.Sprintf("Essential for %s activities", activity),
				})
			}
		}
	}

	// Add destination-specific suggestions
	destinationSuggestions := getDestinationSuggestions(destination, season)
	suggestions = append(suggestions, destinationSuggestions...)

	return suggestions, nil
}

// getSeasonWeatherCategory maps seasons to weather categories
func getSeasonWeatherCategory(season string) string {
	switch strings.ToLower(season) {
	case "summer":
		return "hot"
	case "spring":
		return "warm"
	case "fall", "autumn":
		return "mild"
	case "winter":
		return "cold"
	default:
		return "mild" // default fallback
	}
}

// getDestinationSuggestions provides destination-specific suggestions
func getDestinationSuggestions(destination, season string) []interface{} {
	var suggestions []interface{}

	// Convert destination to lowercase for easier matching
	dest := strings.ToLower(destination)

	// Beach destinations
	if strings.Contains(dest, "beach") || strings.Contains(dest, "hawaii") ||
		strings.Contains(dest, "maldives") || strings.Contains(dest, "bali") {
		suggestions = append(suggestions, map[string]interface{}{
			"item":     "Beach towel",
			"category": "beach",
			"reason":   "Essential for beach destinations",
		})
		suggestions = append(suggestions, map[string]interface{}{
			"item":     "Sunscreen SPF 30+",
			"category": "beach",
			"reason":   "Protection from sun exposure",
		})
	}

	// Mountain destinations
	if strings.Contains(dest, "mountain") || strings.Contains(dest, "ski") ||
		strings.Contains(dest, "alps") || strings.Contains(dest, "rocky") {
		suggestions = append(suggestions, map[string]interface{}{
			"item":     "Warm jacket",
			"category": "mountain",
			"reason":   "Mountains can be cold even in summer",
		})
		suggestions = append(suggestions, map[string]interface{}{
			"item":     "Hiking boots",
			"category": "mountain",
			"reason":   "Sturdy footwear for mountain terrain",
		})
	}

	// City destinations
	if strings.Contains(dest, "city") || strings.Contains(dest, "new york") ||
		strings.Contains(dest, "london") || strings.Contains(dest, "paris") {
		suggestions = append(suggestions, map[string]interface{}{
			"item":     "Comfortable walking shoes",
			"category": "city",
			"reason":   "City exploration involves lots of walking",
		})
		suggestions = append(suggestions, map[string]interface{}{
			"item":     "Power bank",
			"category": "city",
			"reason":   "Keep your phone charged for navigation and photos",
		})
	}

	// Cold weather destinations
	if strings.Contains(dest, "arctic") || strings.Contains(dest, "alaska") ||
		strings.Contains(dest, "iceland") || (strings.Contains(dest, "northern") && season == "winter") {
		suggestions = append(suggestions, map[string]interface{}{
			"item":     "Thermal underwear",
			"category": "cold",
			"reason":   "Essential for cold weather destinations",
		})
		suggestions = append(suggestions, map[string]interface{}{
			"item":     "Hand warmers",
			"category": "cold",
			"reason":   "Keep hands warm in extreme cold",
		})
	}

	return suggestions
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
