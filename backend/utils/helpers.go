package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// GenerateID creates a unique identifier
func GenerateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// LoadJSONFile loads and parses a JSON file
func LoadJSONFile(filepath string, v interface{}) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	err = json.Unmarshal(bytes, v)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	return nil
}

// SaveJSONFile saves data to a JSON file
func SaveJSONFile(filepath string, v interface{}) error {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	err = ioutil.WriteFile(filepath, bytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

// FormatCurrency formats a float as currency
func FormatCurrency(amount float64) string {
	return fmt.Sprintf("$%.2f", amount)
}

// FormatDate formats a time.Time as a readable date
func FormatDate(t time.Time) string {
	return t.Format("January 2, 2006")
}

// FormatDateTime formats a time.Time as a readable date and time
func FormatDateTime(t time.Time) string {
	return t.Format("January 2, 2006 at 3:04 PM")
}

// CalculateDistance calculates the distance between two points using Haversine formula
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth's radius in kilometers

	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad

	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

// ParseFloat safely parses a string to float64
func ParseFloat(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}

// ParseInt safely parses a string to int
func ParseInt(s string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(s))
}

// Contains checks if a slice contains a specific value
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RemoveDuplicates removes duplicate strings from a slice
func RemoveDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// GetSeason determines the season based on month
func GetSeason(month int) string {
	switch month {
	case 12, 1, 2:
		return "winter"
	case 3, 4, 5:
		return "spring"
	case 6, 7, 8:
		return "summer"
	case 9, 10, 11:
		return "fall"
	default:
		return "unknown"
	}
}

// GetWeatherCategory categorizes temperature into weather categories
func GetWeatherCategory(temp float64) string {
	switch {
	case temp >= 25:
		return "hot"
	case temp >= 15:
		return "warm"
	case temp >= 5:
		return "mild"
	case temp >= -5:
		return "cool"
	default:
		return "cold"
	}
}

// CalculateDuration calculates the duration between two dates
func CalculateDuration(start, end time.Time) int {
	duration := end.Sub(start)
	return int(duration.Hours() / 24)
}

// ValidateEmail validates email format (basic)
func ValidateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// SanitizeString removes potentially dangerous characters
func SanitizeString(s string) string {
	// Remove HTML tags
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")

	// Remove script tags
	s = strings.ReplaceAll(s, "script", "")
	s = strings.ReplaceAll(s, "javascript:", "")

	return strings.TrimSpace(s)
}

// CreateDirectory creates a directory if it doesn't exist
func CreateDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

// GetFileSize gets the size of a file in bytes
func GetFileSize(filepath string) (int64, error) {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

// FileExists checks if a file exists
func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

// GetFileExtension gets the file extension
func GetFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

// IsValidImageExtension checks if the file extension is a valid image type
func IsValidImageExtension(ext string) bool {
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}
	return Contains(validExtensions, ext)
}

// LogError logs an error with timestamp
func LogError(message string, err error) {
	log.Printf("[ERROR] %s: %v", message, err)
}

// LogInfo logs an info message with timestamp
func LogInfo(message string) {
	log.Printf("[INFO] %s", message)
}

// LogWarning logs a warning message with timestamp
func LogWarning(message string) {
	log.Printf("[WARNING] %s", message)
}

// RetryWithBackoff retries a function with exponential backoff
func RetryWithBackoff(maxRetries int, fn func() error) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = fn()
		if err == nil {
			return nil
		}

		if i < maxRetries-1 {
			backoff := time.Duration(math.Pow(2, float64(i))) * time.Second
			time.Sleep(backoff)
		}
	}
	return err
}

// TruncateString truncates a string to a maximum length
func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-3] + "..."
}

// CapitalizeFirst capitalizes the first letter of a string
func CapitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

// TitleCase converts a string to title case
func TitleCase(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		words[i] = CapitalizeFirst(word)
	}
	return strings.Join(words, " ")
}
