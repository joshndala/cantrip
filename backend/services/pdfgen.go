package services

//COMPLETED
import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// PDFMetadata represents metadata for a PDF file
type PDFMetadata struct {
	ID            string                 `json:"id"`
	Filename      string                 `json:"filename"`
	Type          string                 `json:"type"` // itinerary, packing, tips
	Size          int64                  `json:"size"`
	CreatedAt     time.Time              `json:"created_at"`
	ExpiresAt     time.Time              `json:"expires_at"`
	DownloadURL   string                 `json:"download_url"`
	ShareURL      string                 `json:"share_url,omitempty"`
	Customization map[string]interface{} `json:"customization,omitempty"`
}

// PDFStatus represents the status of PDF generation
type PDFStatus struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`   // pending, processing, completed, failed
	Progress  int       `json:"progress"` // 0-100
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Error     string    `json:"error,omitempty"`
	URL       string    `json:"url,omitempty"`
}

// Global PDF storage directory
const PDFStorageDir = "data/pdfs"

// Initialize PDF storage
func init() {
	if err := os.MkdirAll(PDFStorageDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create PDF storage directory: %v", err))
	}
}

// GenerateItineraryPDF generates a PDF for an itinerary
func GenerateItineraryPDF(id, format string, includeImages bool, customization map[string]interface{}) (string, error) {
	// Get itinerary data
	itineraryInterface, err := GetItinerary(id)
	if err != nil {
		return "", fmt.Errorf("failed to get itinerary: %w", err)
	}

	// Convert interface{} to map for access
	itineraryData, ok := itineraryInterface.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid itinerary data format")
	}

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Add title
	pdf.Cell(0, 10, "Travel Itinerary")
	pdf.Ln(15)

	// Add itinerary details
	pdf.SetFont("Arial", "B", 12)
	if city, ok := itineraryData["city"].(string); ok {
		pdf.Cell(0, 8, fmt.Sprintf("Destination: %s", city))
		pdf.Ln(10)
	}

	// Handle dates (they might be strings from the interface)
	if startDate, ok := itineraryData["start_date"].(string); ok {
		if endDate, ok := itineraryData["end_date"].(string); ok {
			pdf.Cell(0, 8, fmt.Sprintf("Duration: %s to %s", startDate, endDate))
			pdf.Ln(15)
		}
	}

	// Add daily plans
	if days, ok := itineraryData["days"].([]interface{}); ok {
		for i, dayInterface := range days {
			day, ok := dayInterface.(map[string]interface{})
			if !ok {
				continue
			}

			// Day header
			pdf.SetFont("Arial", "B", 12)
			if dayNum, ok := day["day"].(float64); ok {
				if dateStr, ok := day["date"].(string); ok {
					pdf.Cell(0, 8, fmt.Sprintf("Day %.0f - %s", dayNum, dateStr))
					pdf.Ln(10)
				}
			}

			// Activities
			if activities, ok := day["activities"].([]interface{}); ok && len(activities) > 0 {
				pdf.SetFont("Arial", "B", 10)
				pdf.Cell(0, 6, "Activities:")
				pdf.Ln(8)

				pdf.SetFont("Arial", "", 10)
				for _, activityInterface := range activities {
					if activity, ok := activityInterface.(map[string]interface{}); ok {
						name, _ := activity["name"].(string)
						startTime, _ := activity["start_time"].(string)
						endTime, _ := activity["end_time"].(string)
						pdf.Cell(0, 5, fmt.Sprintf("• %s (%s - %s)", name, startTime, endTime))
						pdf.Ln(6)
					}
				}
				pdf.Ln(5)
			}

			// Meals
			if meals, ok := day["meals"].([]interface{}); ok && len(meals) > 0 {
				pdf.SetFont("Arial", "B", 10)
				pdf.Cell(0, 6, "Meals:")
				pdf.Ln(8)

				pdf.SetFont("Arial", "", 10)
				for _, mealInterface := range meals {
					if meal, ok := mealInterface.(map[string]interface{}); ok {
						mealType, _ := meal["type"].(string)
						name, _ := meal["name"].(string)
						timeStr, _ := meal["time"].(string)
						pdf.Cell(0, 5, fmt.Sprintf("• %s: %s at %s",
							strings.Title(mealType), name, timeStr))
						pdf.Ln(6)
					}
				}
				pdf.Ln(5)
			}

			// Add page break if not last day
			if i < len(days)-1 {
				pdf.AddPage()
			}
		}
	}

	// Save PDF
	filename := fmt.Sprintf("itinerary_%s.pdf", id)
	filepath := filepath.Join(PDFStorageDir, filename)

	if err := pdf.OutputFileAndClose(filepath); err != nil {
		return "", fmt.Errorf("failed to save PDF: %w", err)
	}

	// Get file size
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %w", err)
	}

	// Save metadata
	metadata := PDFMetadata{
		ID:            id,
		Filename:      filename,
		Type:          "itinerary",
		Size:          fileInfo.Size(),
		CreatedAt:     time.Now(),
		ExpiresAt:     time.Now().AddDate(0, 1, 0), // Expires in 1 month
		DownloadURL:   fmt.Sprintf("/api/v1/pdf/download/%s", id),
		Customization: customization,
	}

	if err := savePDFMetadata(metadata); err != nil {
		return "", fmt.Errorf("failed to save metadata: %w", err)
	}

	// Try to upload to GCS if available
	if gcsClient := GetGCSClient(); gcsClient != nil {
		ctx := context.Background()
		objectName := fmt.Sprintf("pdfs/%s", filename)
		if err := gcsClient.UploadFileFromPath(ctx, objectName, filepath); err == nil {
			if signedURL, err := gcsClient.GenerateSignedURL(ctx, objectName, 24*time.Hour); err == nil {
				metadata.DownloadURL = signedURL
				savePDFMetadata(metadata)
			}
		}
	}

	return metadata.DownloadURL, nil
}

// GeneratePackingListPDF generates a PDF for a packing list
func GeneratePackingListPDF(id, format string, includeImages bool, customization map[string]interface{}) (string, error) {
	// Get packing list data
	packingList, err := GetPackingList(id)
	if err != nil {
		return "", fmt.Errorf("failed to get packing list: %w", err)
	}

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Add title
	pdf.Cell(0, 10, "Packing List")
	pdf.Ln(15)

	// Add destination info
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Destination: %s", packingList.Destination))
	pdf.Ln(10)
	pdf.Cell(0, 8, fmt.Sprintf("Total Items: %d", packingList.TotalItems))
	pdf.Ln(15)

	// Add categories and items
	for _, categoryInterface := range packingList.Categories {
		categoryData, _ := json.Marshal(categoryInterface)
		var category PackingCategory
		if err := json.Unmarshal(categoryData, &category); err != nil {
			continue
		}

		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(0, 8, category.Name)
		pdf.Ln(10)

		pdf.SetFont("Arial", "", 10)
		for _, item := range category.Items {
			pdf.Cell(0, 5, fmt.Sprintf("• %s (Qty: %d) - %s",
				item.Name, item.Quantity, item.Reason))
			pdf.Ln(6)
		}
		pdf.Ln(5)
	}

	// Add notes
	if len(packingList.Notes) > 0 {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(0, 8, "Notes:")
		pdf.Ln(10)

		pdf.SetFont("Arial", "", 10)
		for _, note := range packingList.Notes {
			pdf.Cell(0, 5, fmt.Sprintf("• %s", note))
			pdf.Ln(6)
		}
	}

	// Save PDF
	filename := fmt.Sprintf("packing_%s.pdf", id)
	filepath := filepath.Join(PDFStorageDir, filename)

	if err := pdf.OutputFileAndClose(filepath); err != nil {
		return "", fmt.Errorf("failed to save PDF: %w", err)
	}

	// Get file size
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %w", err)
	}

	// Save metadata
	metadata := PDFMetadata{
		ID:            id,
		Filename:      filename,
		Type:          "packing",
		Size:          fileInfo.Size(),
		CreatedAt:     time.Now(),
		ExpiresAt:     time.Now().AddDate(0, 1, 0),
		DownloadURL:   fmt.Sprintf("/api/v1/pdf/download/%s", id),
		Customization: customization,
	}

	if err := savePDFMetadata(metadata); err != nil {
		return "", fmt.Errorf("failed to save metadata: %w", err)
	}

	// Try to upload to GCS if available
	if gcsClient := GetGCSClient(); gcsClient != nil {
		ctx := context.Background()
		objectName := fmt.Sprintf("pdfs/%s", filename)
		if err := gcsClient.UploadFileFromPath(ctx, objectName, filepath); err == nil {
			if signedURL, err := gcsClient.GenerateSignedURL(ctx, objectName, 24*time.Hour); err == nil {
				metadata.DownloadURL = signedURL
				savePDFMetadata(metadata)
			}
		}
	}

	return metadata.DownloadURL, nil
}

// GenerateTipsPDF generates a PDF for travel tips
func GenerateTipsPDF(destination, category string, includeImages bool, customization map[string]interface{}) (string, error) {
	// Get tips data
	tips, err := GetTravelTips(destination, category, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get tips: %w", err)
	}

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Add title
	pdf.Cell(0, 10, fmt.Sprintf("Travel Tips - %s", destination))
	pdf.Ln(15)

	// Add category
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Category: %s", strings.Title(category)))
	pdf.Ln(15)

	// Add tips
	for i, tip := range tips {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(0, 8, fmt.Sprintf("%d. %s", i+1, tip.Title))
		pdf.Ln(10)

		pdf.SetFont("Arial", "", 10)
		pdf.MultiCell(0, 5, tip.Description, "", "", false)
		pdf.Ln(5)

		// Add priority and tags
		pdf.SetFont("Arial", "I", 9)
		pdf.Cell(0, 5, fmt.Sprintf("Priority: %s | Tags: %s",
			tip.Priority, strings.Join(tip.Tags, ", ")))
		pdf.Ln(8)

		// Add examples if available
		if len(tip.Examples) > 0 {
			pdf.SetFont("Arial", "B", 9)
			pdf.Cell(0, 5, "Examples:")
			pdf.Ln(6)

			pdf.SetFont("Arial", "", 9)
			for _, example := range tip.Examples {
				pdf.Cell(0, 4, fmt.Sprintf("• %s", example))
				pdf.Ln(5)
			}
		}

		pdf.Ln(8)

		// Add page break if needed
		if i > 0 && i%3 == 0 {
			pdf.AddPage()
		}
	}

	// Save PDF
	filename := fmt.Sprintf("tips_%s_%s.pdf", strings.ToLower(destination), category)
	filepath := filepath.Join(PDFStorageDir, filename)

	if err := pdf.OutputFileAndClose(filepath); err != nil {
		return "", fmt.Errorf("failed to save PDF: %w", err)
	}

	// Get file size
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %w", err)
	}

	// Generate ID for this PDF
	pdfID := fmt.Sprintf("tips_%s_%s_%d",
		strings.ToLower(destination),
		category,
		time.Now().Unix())

	// Save metadata
	metadata := PDFMetadata{
		ID:            pdfID,
		Filename:      filename,
		Type:          "tips",
		Size:          fileInfo.Size(),
		CreatedAt:     time.Now(),
		ExpiresAt:     time.Now().AddDate(0, 1, 0),
		DownloadURL:   fmt.Sprintf("/api/v1/pdf/download/%s", pdfID),
		Customization: customization,
	}

	if err := savePDFMetadata(metadata); err != nil {
		return "", fmt.Errorf("failed to save metadata: %w", err)
	}

	// Try to upload to GCS if available
	if gcsClient := GetGCSClient(); gcsClient != nil {
		ctx := context.Background()
		objectName := fmt.Sprintf("pdfs/%s", filename)
		if err := gcsClient.UploadFileFromPath(ctx, objectName, filepath); err == nil {
			if signedURL, err := gcsClient.GenerateSignedURL(ctx, objectName, 24*time.Hour); err == nil {
				metadata.DownloadURL = signedURL
				savePDFMetadata(metadata)
			}
		}
	}

	return metadata.DownloadURL, nil
}

// GetPDFMetadata retrieves metadata for a PDF
func GetPDFMetadata(pdfID string) (*PDFMetadata, error) {
	metadata, err := loadPDFMetadata(pdfID)
	if err != nil {
		return nil, fmt.Errorf("failed to load PDF metadata: %w", err)
	}
	return metadata, nil
}

// DownloadPDF downloads a PDF file
func DownloadPDF(id, format string) ([]byte, string, error) {
	// Get metadata
	metadata, err := GetPDFMetadata(id)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get PDF metadata: %w", err)
	}

	// Check if file exists locally
	filepath := filepath.Join(PDFStorageDir, metadata.Filename)
	if _, err := os.Stat(filepath); err == nil {
		// File exists locally
		data, err := os.ReadFile(filepath)
		if err != nil {
			return nil, "", fmt.Errorf("failed to read PDF file: %w", err)
		}
		return data, metadata.Filename, nil
	}

	// Try to download from GCS if available
	if gcsClient := GetGCSClient(); gcsClient != nil {
		ctx := context.Background()
		objectName := fmt.Sprintf("pdfs/%s", metadata.Filename)

		data, err := gcsClient.DownloadFile(ctx, objectName)
		if err == nil {
			// Save locally for future use
			if err := os.WriteFile(filepath, data, 0644); err == nil {
				return data, metadata.Filename, nil
			}
		}
	}

	return nil, "", fmt.Errorf("PDF file not found")
}

// GetPDFStatus checks the status of PDF generation
func GetPDFStatus(id string) (*PDFStatus, error) {
	// For now, assume all PDFs are completed
	// In a real implementation, you'd track generation status
	metadata, err := GetPDFMetadata(id)
	if err != nil {
		return nil, err
	}

	status := &PDFStatus{
		ID:        id,
		Status:    "completed",
		Progress:  100,
		CreatedAt: metadata.CreatedAt,
		UpdatedAt: metadata.CreatedAt,
		URL:       metadata.DownloadURL,
	}

	return status, nil
}

// DeletePDF deletes a PDF file
func DeletePDF(id string) error {
	// Get metadata
	metadata, err := GetPDFMetadata(id)
	if err != nil {
		return fmt.Errorf("failed to get PDF metadata: %w", err)
	}

	// Delete local file
	filepath := filepath.Join(PDFStorageDir, metadata.Filename)
	if err := os.Remove(filepath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete local PDF file: %w", err)
	}

	// Delete from GCS if available
	if gcsClient := GetGCSClient(); gcsClient != nil {
		ctx := context.Background()
		objectName := fmt.Sprintf("pdfs/%s", metadata.Filename)
		if err := gcsClient.DeleteFile(ctx, objectName); err != nil {
			// Log error but don't fail - local deletion is more important
			fmt.Printf("Failed to delete PDF from GCS: %v\n", err)
		}
	}

	// Delete metadata
	if err := deletePDFMetadata(id); err != nil {
		return fmt.Errorf("failed to delete PDF metadata: %w", err)
	}

	return nil
}

// ListUserPDFs lists PDFs for a user (simplified - lists all PDFs)
func ListUserPDFs(userID string) ([]PDFMetadata, error) {
	// In a real implementation, you'd filter by user ID
	// For now, return all PDFs
	return listAllPDFs()
}

// CreateShareableLink creates a shareable link for a PDF
func CreateShareableLink(id, expiryHours string) (string, error) {
	// Get metadata
	metadata, err := GetPDFMetadata(id)
	if err != nil {
		return "", fmt.Errorf("failed to get PDF metadata: %w", err)
	}

	// Parse expiry hours
	hours, err := strconv.Atoi(expiryHours)
	if err != nil {
		hours = 24 // Default to 24 hours
	}

	// Generate shareable link
	shareURL := fmt.Sprintf("/api/v1/pdf/share/%s?expires=%d", id, hours)
	metadata.ShareURL = shareURL
	metadata.ExpiresAt = time.Now().Add(time.Duration(hours) * time.Hour)

	// Save updated metadata
	if err := savePDFMetadata(*metadata); err != nil {
		return "", fmt.Errorf("failed to save shareable link: %w", err)
	}

	return shareURL, nil
}

// Helper functions for metadata management
func savePDFMetadata(metadata PDFMetadata) error {
	metadataFile := filepath.Join(PDFStorageDir, fmt.Sprintf("%s_metadata.json", metadata.ID))
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(metadataFile, data, 0644)
}

func loadPDFMetadata(id string) (*PDFMetadata, error) {
	metadataFile := filepath.Join(PDFStorageDir, fmt.Sprintf("%s_metadata.json", id))
	data, err := os.ReadFile(metadataFile)
	if err != nil {
		return nil, err
	}

	var metadata PDFMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

func deletePDFMetadata(id string) error {
	metadataFile := filepath.Join(PDFStorageDir, fmt.Sprintf("%s_metadata.json", id))
	return os.Remove(metadataFile)
}

func listAllPDFs() ([]PDFMetadata, error) {
	files, err := os.ReadDir(PDFStorageDir)
	if err != nil {
		return nil, err
	}

	var pdfs []PDFMetadata
	for _, file := range files {
		if strings.HasSuffix(file.Name(), "_metadata.json") {
			id := strings.TrimSuffix(file.Name(), "_metadata.json")
			if metadata, err := loadPDFMetadata(id); err == nil {
				pdfs = append(pdfs, *metadata)
			}
		}
	}

	return pdfs, nil
}
