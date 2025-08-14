package services

// GetCachedSuggestions retrieves cached suggestions for a mood and city
func GetCachedSuggestions(mood, city string) ([]TripSuggestion, error) {
	// TODO: Implement actual caching logic
	return []TripSuggestion{
		{
			Title:         "Cached Trip",
			Description:   "A cached trip suggestion for " + city,
			Activities:    []string{"Cached activity 1", "Cached activity 2"},
			EstimatedCost: 500.0,
			Duration:      7,
			Tags:          []string{mood, city},
		},
	}, nil
}
