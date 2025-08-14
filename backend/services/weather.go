package services

type WeatherInfo struct {
	Temperature float64 `json:"temperature"`
	Condition   string  `json:"condition"`
	Humidity    int     `json:"humidity"`
	WindSpeed   float64 `json:"wind_speed"`
}

// GetWeather retrieves weather information for a city
func GetWeather(city string) (WeatherInfo, error) {
	// TODO: Implement actual weather API call
	return WeatherInfo{
		Temperature: 22.0,
		Condition:   "Sunny",
		Humidity:    65,
		WindSpeed:   10.0,
	}, nil
}
