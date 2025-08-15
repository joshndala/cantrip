package services

//COMPLETED
import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// TipsData represents the structure of tips.json
type TipsData struct {
	GeneralCanada map[string]interface{}            `json:"general_canada"`
	Cities        map[string]map[string]interface{} `json:"-"`
}

// Tip represents a travel tip
type Tip struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Priority    string   `json:"priority"`
	Tags        []string `json:"tags"`
	Examples    []string `json:"examples,omitempty"`
}

// Emergency represents emergency information
type Emergency struct {
	Police      string `json:"police"`
	Ambulance   string `json:"ambulance"`
	Fire        string `json:"fire"`
	TouristHelp string `json:"tourist_help"`
}

// Language represents language information
type Language struct {
	Primary       string        `json:"primary"`
	Secondary     []string      `json:"secondary"`
	CommonPhrases []interface{} `json:"common_phrases"`
}

// Global tips data cache
var tipsData *TipsData

// loadTipsData loads tips data from JSON file
func loadTipsData() (*TipsData, error) {
	if tipsData != nil {
		return tipsData, nil
	}

	data, err := os.ReadFile("data/tips.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read tips.json: %w", err)
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tips.json: %w", err)
	}

	// Extract general_canada data
	generalCanada, ok := rawData["general_canada"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("general_canada data not found or invalid")
	}

	// Extract city-specific data
	cities := make(map[string]map[string]interface{})
	for key, value := range rawData {
		if key != "general_canada" {
			if cityData, ok := value.(map[string]interface{}); ok {
				cities[key] = cityData
			}
		}
	}

	tipsData = &TipsData{
		GeneralCanada: generalCanada,
		Cities:        cities,
	}

	return tipsData, nil
}

// mergeTips merges general and city-specific tips with deduplication
func mergeTips(generalTips []Tip, cityTips []Tip) []Tip {
	// Start with general tips
	merged := make([]Tip, len(generalTips))
	copy(merged, generalTips)

	// Create a map for quick title lookup
	titleMap := make(map[string]int)
	for i, tip := range merged {
		titleMap[strings.ToLower(tip.Title)] = i
	}

	// Add city tips, replacing duplicates by title
	for _, cityTip := range cityTips {
		titleLower := strings.ToLower(cityTip.Title)
		if existingIndex, exists := titleMap[titleLower]; exists {
			// Replace existing tip with city-specific one
			merged[existingIndex] = cityTip
		} else {
			// Add new tip
			merged = append(merged, cityTip)
			titleMap[titleLower] = len(merged) - 1
		}
	}

	return merged
}

// extractTipsFromInterface extracts tips from interface{} data
func extractTipsFromInterface(data interface{}, category string) ([]Tip, error) {
	var tips []Tip

	if categoryData, ok := data.(map[string]interface{}); ok {
		if tipsArray, ok := categoryData[category].([]interface{}); ok {
			for _, tipInterface := range tipsArray {
				if tipMap, ok := tipInterface.(map[string]interface{}); ok {
					tip := Tip{
						Category: category,
					}

					if title, ok := tipMap["title"].(string); ok {
						tip.Title = title
					}
					if description, ok := tipMap["description"].(string); ok {
						tip.Description = description
					}
					if priority, ok := tipMap["priority"].(string); ok {
						tip.Priority = priority
					}
					if tags, ok := tipMap["tags"].([]interface{}); ok {
						for _, tag := range tags {
							if tagStr, ok := tag.(string); ok {
								tip.Tags = append(tip.Tags, tagStr)
							}
						}
					}
					if examples, ok := tipMap["examples"].([]interface{}); ok {
						for _, example := range examples {
							if exampleStr, ok := example.(string); ok {
								tip.Examples = append(tip.Examples, exampleStr)
							}
						}
					}

					tips = append(tips, tip)
				}
			}
		}
	}

	return tips, nil
}

// GetTravelTips gets travel tips for a destination
func GetTravelTips(destination, category string, topics []string) ([]Tip, error) {
	data, err := loadTipsData()
	if err != nil {
		return nil, err
	}

	// Get general Canada tips for the category
	generalTips, err := extractTipsFromInterface(data.GeneralCanada, category)
	if err != nil {
		return nil, err
	}

	// Get city-specific tips if available
	var cityTips []Tip
	if cityData, exists := data.Cities[destination]; exists {
		cityTips, err = extractTipsFromInterface(cityData, category)
		if err != nil {
			return nil, err
		}
	}

	// Merge tips with deduplication
	mergedTips := mergeTips(generalTips, cityTips)

	// Filter by topics if provided
	if len(topics) > 0 {
		filteredTips := []Tip{}
		for _, tip := range mergedTips {
			for _, topic := range topics {
				for _, tag := range tip.Tags {
					if strings.Contains(strings.ToLower(tag), strings.ToLower(topic)) {
						filteredTips = append(filteredTips, tip)
						break
					}
				}
			}
		}
		mergedTips = filteredTips
	}

	return mergedTips, nil
}

// GetEmergencyInfo gets emergency information for a destination
func GetEmergencyInfo(destination string) (Emergency, error) {
	data, err := loadTipsData()
	if err != nil {
		return Emergency{}, err
	}

	// Get general Canada emergency info
	generalEmergency := Emergency{
		Police:      "911",
		Ambulance:   "911",
		Fire:        "911",
		TouristHelp: "1-800-TOURIST",
	}

	if emergencyData, ok := data.GeneralCanada["emergency"].(map[string]interface{}); ok {
		if general, ok := emergencyData["general"].(string); ok {
			generalEmergency.Police = general
			generalEmergency.Ambulance = general
			generalEmergency.Fire = general
		}
	}

	// City-specific emergency info could be added here if available
	// For now, return general Canada emergency info

	return generalEmergency, nil
}

// GetLanguageInfo gets language information for a destination
func GetLanguageInfo(destination string) (Language, error) {
	data, err := loadTipsData()
	if err != nil {
		return Language{}, err
	}

	// Get general Canada language info
	language := Language{
		Primary:   "English",
		Secondary: []string{"French"},
		CommonPhrases: []interface{}{
			map[string]string{
				"english": "Hello",
				"local":   "Bonjour",
			},
		},
	}

	if languageData, ok := data.GeneralCanada["language"].(map[string]interface{}); ok {
		if primary, ok := languageData["primary"].([]interface{}); ok {
			language.Primary = primary[0].(string)
			for i := 1; i < len(primary); i++ {
				if lang, ok := primary[i].(string); ok {
					language.Secondary = append(language.Secondary, lang)
				}
			}
		}
		if phrases, ok := languageData["common_phrases"].([]interface{}); ok {
			language.CommonPhrases = phrases
		}
	}

	// Get city-specific language info if available
	if cityData, exists := data.Cities[destination]; exists {
		if cityLanguageData, ok := cityData["language"].(map[string]interface{}); ok {
			if notes, ok := cityLanguageData["notes"].(string); ok {
				// Add city-specific language notes
				language.CommonPhrases = append(language.CommonPhrases, map[string]string{
					"note": notes,
				})
			}
		}
	}

	return language, nil
}

// GetCulturalTips gets cultural tips for a destination
func GetCulturalTips(destination string) ([]Tip, error) {
	return GetTravelTips(destination, "cultural", nil)
}

// GetTippingGuide gets tipping guide for a destination
func GetTippingGuide(destination string) (map[string]interface{}, error) {
	tips, err := GetTravelTips(destination, "tipping", nil)
	if err != nil {
		return nil, err
	}

	// Convert tips to a more structured format for tipping guide
	guide := make(map[string]interface{})
	for _, tip := range tips {
		guide[tip.Title] = map[string]interface{}{
			"description": tip.Description,
			"priority":    tip.Priority,
			"tags":        tip.Tags,
		}
	}

	return guide, nil
}

// GetSafetyTips gets safety tips for a destination
func GetSafetyTips(destination string) ([]Tip, error) {
	return GetTravelTips(destination, "safety", nil)
}

// GetLocalCustoms gets local customs for a destination
func GetLocalCustoms(destination string) ([]Tip, error) {
	return GetTravelTips(destination, "customs", nil)
}
