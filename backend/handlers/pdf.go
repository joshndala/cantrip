package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joshndala/cantrip/services"
)

type PDFRequest struct {
	Type          string                 `json:"type" binding:"required"` // "itinerary", "packing", "tips"
	ID            string                 `json:"id" binding:"required"`
	Format        string                 `json:"format"` // "pdf", "html"
	IncludeImages bool                   `json:"include_images"`
	Customization map[string]interface{} `json:"customization"`
}

type PDFResponse struct {
	URL         string `json:"url"`
	Filename    string `json:"filename"`
	Size        int64  `json:"size"` // in bytes
	ExpiresAt   string `json:"expires_at"`
	DownloadURL string `json:"download_url"`
}

// GeneratePDFHandler creates downloadable PDFs for itineraries, packing lists, etc.
func GeneratePDFHandler(c *gin.Context) {
	var req PDFRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var pdfURL string
	var err error

	switch req.Type {
	case "itinerary":
		pdfURL, err = services.GenerateItineraryPDF(req.ID, req.Format, req.IncludeImages, req.Customization)
	case "packing":
		pdfURL, err = services.GeneratePackingListPDF(req.ID, req.Format, req.IncludeImages, req.Customization)
	case "tips":
		pdfURL, err = services.GenerateTipsPDF(req.ID, req.Format, req.IncludeImages, req.Customization)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PDF type"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF: " + err.Error()})
		return
	}

	// Get file metadata
	metadata, err := services.GetPDFMetadata(pdfURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get PDF metadata"})
		return
	}

	response := PDFResponse{
		URL:         pdfURL,
		Filename:    metadata["filename"].(string),
		Size:        metadata["size"].(int64),
		ExpiresAt:   metadata["expires_at"].(string),
		DownloadURL: metadata["download_url"].(string),
	}

	c.JSON(http.StatusOK, response)
}

// DownloadPDFHandler serves PDF files for download
func DownloadPDFHandler(c *gin.Context) {
	id := c.Param("id")
	format := c.Query("format")
	if format == "" {
		format = "pdf"
	}

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File ID is required"})
		return
	}

	// Get file from GCS
	fileData, filename, err := services.DownloadPDF(id, format)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Set headers for file download
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Length", strconv.Itoa(len(fileData)))

	c.Data(http.StatusOK, "application/pdf", fileData)
}

// GetPDFStatusHandler checks the status of PDF generation
func GetPDFStatusHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File ID is required"})
		return
	}

	status, err := services.GetPDFStatus(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.JSON(http.StatusOK, status)
}

// DeletePDFHandler deletes a generated PDF
func DeletePDFHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File ID is required"})
		return
	}

	err := services.DeletePDF(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete PDF"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PDF deleted successfully"})
}

// ListPDFsHandler lists all PDFs for a user
func ListPDFsHandler(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	pdfs, err := services.ListUserPDFs(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list PDFs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"pdfs":    pdfs,
	})
}

// SharePDFHandler generates a shareable link for a PDF
func SharePDFHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File ID is required"})
		return
	}

	expiryHours := c.Query("expiry_hours")
	if expiryHours == "" {
		expiryHours = "24" // Default 24 hours
	}

	shareURL, err := services.CreateShareableLink(id, expiryHours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create shareable link"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"share_url":  shareURL,
		"expires_in": expiryHours + " hours",
	})
}
