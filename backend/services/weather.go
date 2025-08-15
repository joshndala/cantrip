package services

//COMPLETED
import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

// WeatherInfo represents weather information for a location
type WeatherInfo struct {
	Temperature float64 `json:"temperature"`
	Condition   string  `json:"condition"`
	Humidity    int     `json:"humidity"`
	WindSpeed   float64 `json:"wind_speed"`
}

// WeatherForecast represents a weather forecast for a specific date
type WeatherForecast struct {
	Date          string  `json:"date"`
	HighTemp      float64 `json:"high_temp"`
	LowTemp       float64 `json:"low_temp"`
	Condition     string  `json:"condition"`
	Humidity      int     `json:"humidity"`
	WindSpeed     float64 `json:"wind_speed"`
	Precipitation float64 `json:"precipitation"`
}

// WeatherForecastResponse represents the response from OpenWeatherMap forecast API
type WeatherForecastResponse struct {
	List []ForecastItem `json:"list"`
	City CityInfo       `json:"city"`
}

// ForecastItem represents a 3-hour forecast period
type ForecastItem struct {
	Dt   int64 `json:"dt"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	POP  *float64           `json:"pop,omitempty"`
	Rain map[string]float64 `json:"rain,omitempty"`
	Snow map[string]float64 `json:"snow,omitempty"`
}

// CityInfo represents city information from the API
type CityInfo struct {
	Name     string `json:"name"`
	Country  string `json:"country"`
	Timezone int    `json:"timezone"`
}

// CityMetadata represents the structure of city_metadata.json
type CityMetadata struct {
	Cities []City `json:"cities"`
}

// City represents a single city in the metadata
type City struct {
	Name          string            `json:"name"`
	Province      string            `json:"province"`
	Country       string            `json:"country"`
	Coordinates   Coordinates       `json:"coordinates"`
	Timezone      string            `json:"timezone"`
	Population    int               `json:"population"`
	Description   string            `json:"description"`
	Seasons       map[string]Season `json:"seasons"`
	Attractions   []string          `json:"attractions"`
	Neighborhoods []string          `json:"neighborhoods"`
}

// Coordinates represents latitude and longitude
type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Season represents seasonal weather data
type Season struct {
	Months     []string `json:"months"`
	AvgTemp    float64  `json:"avg_temp"`
	Activities []string `json:"activities"`
}

// GetWeather retrieves weather information for a city
func GetWeather(city string) (WeatherInfo, error) {
	// First, try to get weather from a real API (if configured)
	if weather, err := getWeatherFromAPI(city); err == nil {
		return weather, nil
	}

	// Fallback to using city metadata for seasonal weather
	return getWeatherFromMetadata(city)
}

// GetWeatherForecast retrieves weather forecast for a city and trip dates
func GetWeatherForecast(city string, startDate, endDate string) ([]WeatherForecast, error) {
	// Parse trip dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}

	// Calculate days from today using city timezone (will be updated when we get API response)
	today := time.Now().Truncate(24 * time.Hour)
	daysFromToday := int(start.Sub(today).Hours() / 24)

	// If trip is within 5 days, get forecast from API
	if daysFromToday <= 5 {
		// Try to get real forecast for the entire trip or first 5 days
		realForecast, err := getForecastFromAPI(city, start, end)
		if err == nil {
			// If trip extends beyond 5 days, add seasonal data for remaining days
			if end.After(start.AddDate(0, 0, 5)) {
				seasonalForecast, err := getSeasonalForecast(city, start.AddDate(0, 0, 6), end)
				if err == nil {
					return append(realForecast, seasonalForecast...), nil
				}
			}
			return realForecast, nil
		}
	}

	// If trip is beyond 5 days or API failed, use seasonal data
	return getSeasonalForecast(city, start, end)
}

// GetWeatherForecastWithNotes retrieves weather forecast with helpful notes
func GetWeatherForecastWithNotes(city string, startDate, endDate string) ([]WeatherForecast, []string, error) {
	forecasts, err := GetWeatherForecast(city, startDate, endDate)
	if err != nil {
		return nil, nil, err
	}

	// Parse dates for note generation
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	today := time.Now().Truncate(24 * time.Hour)
	daysFromToday := int(start.Sub(today).Hours() / 24)

	var notes []string

	// Add note for trips beyond 5 days
	if daysFromToday > 5 {
		notes = append(notes, getWeatherNoteForLongTermTrip(daysFromToday))
	} else if end.After(start.AddDate(0, 0, 5)) {
		// Hybrid forecast: real data for first 5 days, seasonal for rest
		notes = append(notes, "Forecast combines real-time data for the first 5 days with seasonal averages for the remaining days. Check closer to your trip for updates on the later dates.")
	}

	// Add seasonal notes
	seasonalNotes := getSeasonalWeatherNotes(city, start, end)
	notes = append(notes, seasonalNotes...)

	return forecasts, notes, nil
}

// getWeatherNoteForLongTermTrip generates a note for trips beyond 5 days
func getWeatherNoteForLongTermTrip(daysFromToday int) string {
	if daysFromToday <= 7 {
		return "Weather forecast is based on seasonal averages. Check closer to your trip date for more accurate predictions."
	} else if daysFromToday <= 14 {
		return "Long-term weather forecast uses seasonal data. Consider checking weather updates 1-2 weeks before your trip."
	} else {
		return "Extended forecast uses historical seasonal data. Weather patterns can vary significantly, so check closer to your travel dates."
	}
}

// resolveCityCoordinates resolves city name to coordinates for more accurate API calls
// This helps with city disambiguation (e.g., "London" could be UK or Ontario)
func resolveCityCoordinates(city string) (float64, float64, error) {
	// For now, return error to use city name directly
	// In the future, this could call OpenWeatherMap's Geocoding API
	// Example: http://api.openweathermap.org/geo/1.0/direct?q=London&limit=5&appid={API_KEY}
	return 0, 0, fmt.Errorf("coordinates not resolved, using city name")
}

// getForecastByCoordinates gets forecast using lat/lon instead of city name
func getForecastByCoordinates(lat, lon float64, start, end time.Time) ([]WeatherForecast, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("no weather API key configured")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Use coordinates for more precise location
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?lat=%.4f&lon=%.4f&appid=%s&units=metric", lat, lon, apiKey)

	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch forecast: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("forecast API returned status: %d", resp.StatusCode)
	}

	var forecastResp WeatherForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecastResp); err != nil {
		return nil, fmt.Errorf("failed to decode forecast response: %w", err)
	}

	// Convert trip dates to city's timezone for proper comparison
	tz := time.FixedZone("city", forecastResp.City.Timezone)
	startLocal := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, tz)
	endLocal := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, tz)

	return aggregateForecastData(forecastResp, startLocal, endLocal)
}

// getSeasonalWeatherNotes generates helpful notes based on seasonal weather patterns
func getSeasonalWeatherNotes(city string, start, end time.Time) []string {
	var notes []string

	// Get seasons for the trip period
	startSeason := getSeasonForDate(start)
	endSeason := getSeasonForDate(end)

	if startSeason != endSeason {
		notes = append(notes, fmt.Sprintf("Your trip spans %s to %s seasons. Pack versatile clothing for changing weather.", startSeason, endSeason))
	} else {
		notes = append(notes, fmt.Sprintf("Your trip is during %s. Pack accordingly for typical %s weather in %s.", startSeason, startSeason, city))
	}

	// Add season-specific advice
	switch startSeason {
	case "winter":
		notes = append(notes, "Winter travel tip: Pack layers and warm accessories. Weather can be unpredictable with potential snow or rain.")
	case "spring":
		notes = append(notes, "Spring travel tip: Weather can be variable. Pack layers and be prepared for both warm and cool days.")
	case "summer":
		notes = append(notes, "Summer travel tip: Expect warm weather. Don't forget sun protection and lightweight clothing.")
	case "fall":
		notes = append(notes, "Fall travel tip: Temperatures can drop significantly. Pack layers and warm clothing for cooler evenings.")
	}

	return notes
}

// getWeatherFromAPI attempts to get weather from a real weather API
func getWeatherFromAPI(city string) (WeatherInfo, error) {
	// Check if API key is configured (you can set this as an environment variable)
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return WeatherInfo{}, fmt.Errorf("no weather API key configured")
	}

	// Example using OpenWeatherMap API
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return WeatherInfo{}, fmt.Errorf("failed to fetch weather from API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return WeatherInfo{}, fmt.Errorf("weather API returned status: %d", resp.StatusCode)
	}

	// Parse the API response (simplified - you'd need to implement full parsing)
	var apiResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return WeatherInfo{}, fmt.Errorf("failed to decode API response: %w", err)
	}

	// Extract weather data from API response
	// This is a simplified example - actual implementation would depend on the API
	if main, ok := apiResponse["main"].(map[string]interface{}); ok {
		if temp, ok := main["temp"].(float64); ok {
			if weather, ok := apiResponse["weather"].([]interface{}); ok && len(weather) > 0 {
				if weatherObj, ok := weather[0].(map[string]interface{}); ok {
					if condition, ok := weatherObj["main"].(string); ok {
						return WeatherInfo{
							Temperature: temp,
							Condition:   condition,
							Humidity:    int(main["humidity"].(float64)),
							WindSpeed:   getWindSpeedFromAPI(apiResponse),
						}, nil
					}
				}
			}
		}
	}

	return WeatherInfo{}, fmt.Errorf("failed to parse API response")
}

// getWindSpeedFromAPI extracts wind speed from API response
func getWindSpeedFromAPI(response map[string]interface{}) float64 {
	if wind, ok := response["wind"].(map[string]interface{}); ok {
		if speed, ok := wind["speed"].(float64); ok {
			return speed
		}
	}
	return 0.0
}

// getWeatherFromMetadata gets weather information from city metadata
func getWeatherFromMetadata(city string) (WeatherInfo, error) {
	// Load city metadata
	metadata, err := loadCityMetadata()
	if err != nil {
		return WeatherInfo{}, fmt.Errorf("failed to load city metadata: %w", err)
	}

	// Find the city in metadata
	cityData, err := findCity(metadata, city)
	if err != nil {
		return WeatherInfo{}, err
	}

	// Get current season
	currentSeason := getCurrentSeason()

	// Get seasonal weather data
	seasonData, exists := cityData.Seasons[currentSeason]
	if !exists {
		return WeatherInfo{}, fmt.Errorf("no seasonal data available for %s in %s", currentSeason, city)
	}

	// Generate realistic weather based on seasonal averages
	weather := generateSeasonalWeather(seasonData, currentSeason)

	return weather, nil
}

// loadCityMetadata loads the city metadata from JSON file
func loadCityMetadata() (*CityMetadata, error) {
	data, err := os.ReadFile("data/city_metadata.json")
	if err != nil {
		return nil, err
	}

	var metadata CityMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// findCity finds a city in the metadata by name (case-insensitive)
func findCity(metadata *CityMetadata, cityName string) (*City, error) {
	cityNameLower := strings.ToLower(cityName)

	for _, city := range metadata.Cities {
		if strings.ToLower(city.Name) == cityNameLower {
			return &city, nil
		}
	}

	return nil, fmt.Errorf("city '%s' not found in metadata", cityName)
}

// getCurrentSeason determines the current season based on the current month
func getCurrentSeason() string {
	month := time.Now().Month()

	switch month {
	case time.March, time.April, time.May:
		return "spring"
	case time.June, time.July, time.August:
		return "summer"
	case time.September, time.October, time.November:
		return "fall"
	default:
		return "winter"
	}
}

// generateSeasonalWeather creates realistic weather data based on seasonal averages
func generateSeasonalWeather(season Season, seasonName string) WeatherInfo {
	// Use the average temperature as a base
	baseTemp := season.AvgTemp

	// Add some realistic variation (±5 degrees)
	variation := (rand.Float64() - 0.5) * 10
	temperature := baseTemp + variation

	// Determine weather condition based on season and temperature
	condition := getWeatherCondition(seasonName, temperature)

	// Generate realistic humidity based on condition
	humidity := getHumidityForCondition(condition)

	// Generate realistic wind speed
	windSpeed := getWindSpeedForSeason(seasonName)

	return WeatherInfo{
		Temperature: temperature,
		Condition:   condition,
		Humidity:    humidity,
		WindSpeed:   windSpeed,
	}
}

// getWeatherCondition determines weather condition based on season and temperature
func getWeatherCondition(season string, temperature float64) string {

	switch season {
	case "summer":
		if temperature > 25 {
			return "Sunny"
		} else if temperature > 20 {
			return "Partly Cloudy"
		} else {
			return "Cloudy"
		}
	case "spring":
		if temperature > 15 {
			return "Partly Cloudy"
		} else {
			return "Cloudy"
		}
	case "fall":
		if temperature > 10 {
			return "Partly Cloudy"
		} else {
			return "Cloudy"
		}
	case "winter":
		if temperature < 0 {
			return "Snowy"
		} else if temperature < 5 {
			return "Rainy"
		} else {
			return "Cloudy"
		}
	default:
		return "Partly Cloudy"
	}
}

// getHumidityForCondition returns realistic humidity based on weather condition
func getHumidityForCondition(condition string) int {
	switch condition {
	case "Sunny":
		return rand.Intn(30) + 30 // 30-60%
	case "Partly Cloudy":
		return rand.Intn(20) + 50 // 50-70%
	case "Cloudy":
		return rand.Intn(20) + 60 // 60-80%
	case "Rainy":
		return rand.Intn(20) + 70 // 70-90%
	case "Snowy":
		return rand.Intn(20) + 60 // 60-80%
	default:
		return rand.Intn(30) + 50 // 50-80%
	}
}

// getWindSpeedForSeason returns realistic wind speed based on season
func getWindSpeedForSeason(season string) float64 {
	switch season {
	case "winter":
		return rand.Float64()*15 + 5 // 5-20 km/h
	case "spring":
		return rand.Float64()*20 + 10 // 10-30 km/h
	case "summer":
		return rand.Float64()*10 + 5 // 5-15 km/h
	case "fall":
		return rand.Float64()*15 + 8 // 8-23 km/h
	default:
		return rand.Float64()*10 + 5 // 5-15 km/h
	}
}

// getForecastFromAPI gets weather forecast from OpenWeatherMap API
func getForecastFromAPI(city string, start, end time.Time) ([]WeatherForecast, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("no weather API key configured")
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Get forecast data (5 days, 3-hour intervals)
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?q=%s&appid=%s&units=metric", city, apiKey)

	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch forecast: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("forecast API returned status: %d", resp.StatusCode)
	}

	var forecastResp WeatherForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecastResp); err != nil {
		return nil, fmt.Errorf("failed to decode forecast response: %w", err)
	}

	// Convert trip dates to city's timezone for proper comparison
	tz := time.FixedZone("city", forecastResp.City.Timezone)
	startLocal := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, tz)
	endLocal := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, tz)

	// Aggregate 3-hour forecasts into daily forecasts
	return aggregateForecastData(forecastResp, startLocal, endLocal)
}

// aggregateForecastData aggregates 3-hour forecasts into daily forecasts
func aggregateForecastData(forecastResp WeatherForecastResponse, start, end time.Time) ([]WeatherForecast, error) {
	// Group forecasts by date
	dailyForecasts := make(map[string][]ForecastItem)

	for _, item := range forecastResp.List {
		// Convert timestamp to city's timezone
		itemTime := time.Unix(item.Dt, 0).UTC()
		if forecastResp.City.Timezone != 0 {
			itemTime = itemTime.Add(time.Duration(forecastResp.City.Timezone) * time.Second)
		}

		dateStr := itemTime.Format("2006-01-02")

		// Only include forecasts within our trip dates (compare by date, not timestamp)
		if itemTime.Before(start) || itemTime.After(end) {
			continue
		}

		dailyForecasts[dateStr] = append(dailyForecasts[dateStr], item)
	}

	var forecasts []WeatherForecast

	// Generate daily forecast for each date
	for dateStr, items := range dailyForecasts {
		if len(items) == 0 {
			continue
		}

		// Calculate daily aggregates
		var temps []float64
		var humidities []int
		var windSpeeds []float64
		var precipitations []float64
		var conditions []string
		var pops []float64

		for _, item := range items {
			temps = append(temps, item.Main.Temp)
			humidities = append(humidities, item.Main.Humidity)
			windSpeeds = append(windSpeeds, item.Wind.Speed)

			// Handle precipitation (rain and snow)
			if item.Rain != nil {
				if rain3h, exists := item.Rain["3h"]; exists {
					precipitations = append(precipitations, rain3h)
				}
			}
			if item.Snow != nil {
				if snow3h, exists := item.Snow["3h"]; exists {
					precipitations = append(precipitations, snow3h)
				}
			}

			// Handle probability of precipitation
			if item.POP != nil {
				pops = append(pops, *item.POP)
			}

			if len(item.Weather) > 0 {
				conditions = append(conditions, item.Weather[0].Main)
			}
		}

		// Calculate daily statistics
		highTemp := maxFloat64(temps)
		lowTemp := minFloat64(temps)
		avgHumidity := averageInt(humidities)
		avgWindSpeed := averageFloat64(windSpeeds)
		totalPrecipitation := sumFloat64(precipitations)
		mostCommonCondition := mostCommonString(conditions)
		// Note: avgPOP is calculated but not currently used in WeatherForecast struct
		// Could be added to the struct if needed for future features

		forecasts = append(forecasts, WeatherForecast{
			Date:          dateStr,
			HighTemp:      highTemp,
			LowTemp:       lowTemp,
			Condition:     mostCommonCondition,
			Humidity:      avgHumidity,
			WindSpeed:     avgWindSpeed,
			Precipitation: totalPrecipitation,
		})
	}

	return forecasts, nil
}

// getSeasonalForecast generates forecast based on seasonal data
func getSeasonalForecast(city string, start, end time.Time) ([]WeatherForecast, error) {
	// Load city metadata
	metadata, err := loadCityMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to load city metadata: %w", err)
	}

	// Find the city
	cityData, err := findCity(metadata, city)
	if err != nil {
		return nil, err
	}

	var forecasts []WeatherForecast

	// Generate forecast for each day
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		season := getSeasonForDate(d)
		seasonData, exists := cityData.Seasons[season]
		if !exists {
			continue
		}

		// Generate realistic daily weather
		weather := generateSeasonalWeather(seasonData, season)

		// Add some variation for high/low temps
		variation := (rand.Float64() - 0.5) * 8 // ±4 degrees
		highTemp := weather.Temperature + math.Abs(variation)
		lowTemp := weather.Temperature - math.Abs(variation)

		forecasts = append(forecasts, WeatherForecast{
			Date:          d.Format("2006-01-02"),
			HighTemp:      highTemp,
			LowTemp:       lowTemp,
			Condition:     weather.Condition,
			Humidity:      weather.Humidity,
			WindSpeed:     weather.WindSpeed,
			Precipitation: 0, // Seasonal forecasts don't include precipitation
		})
	}

	return forecasts, nil
}

// getSeasonForDate determines the season for a specific date
func getSeasonForDate(date time.Time) string {
	month := date.Month()

	switch month {
	case time.March, time.April, time.May:
		return "spring"
	case time.June, time.July, time.August:
		return "summer"
	case time.September, time.October, time.November:
		return "fall"
	default:
		return "winter"
	}
}

// Helper functions for aggregating data
func maxFloat64(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}

func minFloat64(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}

func averageInt(values []int) int {
	if len(values) == 0 {
		return 0
	}
	sum := 0
	for _, v := range values {
		sum += v
	}
	return sum / len(values)
}

func averageFloat64(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func sumFloat64(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum
}

func mostCommonString(values []string) string {
	if len(values) == 0 {
		return "Unknown"
	}

	counts := make(map[string]int)
	for _, v := range values {
		counts[v]++
	}

	mostCommon := values[0]
	maxCount := counts[values[0]]

	for value, count := range counts {
		if count > maxCount {
			mostCommon = value
			maxCount = count
		}
	}

	return mostCommon
}
