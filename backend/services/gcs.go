package services

//COMPLETED
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// GCSConfig holds Google Cloud Storage configuration
type GCSConfig struct {
	ProjectID       string
	BucketName      string
	CredentialsFile string
}

// GCSClient wraps the Google Cloud Storage client
type GCSClient struct {
	client     *storage.Client
	bucket     *storage.BucketHandle
	projectID  string
	bucketName string
}

// FileInfo represents information about a stored file
type FileInfo struct {
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	URL         string    `json:"url"`
	Bucket      string    `json:"bucket"`
}

// NewGCSClient creates a new Google Cloud Storage client
func NewGCSClient(config GCSConfig) (*GCSClient, error) {
	ctx := context.Background()

	var client *storage.Client
	var err error

	// Use credentials file if provided, otherwise use default credentials
	if config.CredentialsFile != "" {
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(config.CredentialsFile))
	} else {
		client, err = storage.NewClient(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}

	bucket := client.Bucket(config.BucketName)

	return &GCSClient{
		client:     client,
		bucket:     bucket,
		projectID:  config.ProjectID,
		bucketName: config.BucketName,
	}, nil
}

// UploadFile uploads a file to Google Cloud Storage
func (g *GCSClient) UploadFile(ctx context.Context, objectName string, data []byte, contentType string) error {
	obj := g.bucket.Object(objectName)
	writer := obj.NewWriter(ctx)

	// Set content type
	writer.ContentType = contentType

	// Write data
	if _, err := writer.Write(data); err != nil {
		return fmt.Errorf("failed to write data to GCS: %w", err)
	}

	// Close writer
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close GCS writer: %w", err)
	}

	return nil
}

// UploadFileFromPath uploads a file from local path to Google Cloud Storage
func (g *GCSClient) UploadFileFromPath(ctx context.Context, localPath, objectName string) error {
	// Read local file
	data, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("failed to read local file: %w", err)
	}

	// Determine content type based on file extension
	contentType := getContentTypeFromExtension(filepath.Ext(localPath))

	// Upload to GCS
	return g.UploadFile(ctx, objectName, data, contentType)
}

// DownloadFile downloads a file from Google Cloud Storage
func (g *GCSClient) DownloadFile(ctx context.Context, objectName string) ([]byte, error) {
	obj := g.bucket.Object(objectName)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS reader: %w", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read data from GCS: %w", err)
	}

	return data, nil
}

// DownloadFileToPath downloads a file from GCS to a local path
func (g *GCSClient) DownloadFileToPath(ctx context.Context, objectName, localPath string) error {
	data, err := g.DownloadFile(ctx, objectName)
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write to local file
	if err := os.WriteFile(localPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write local file: %w", err)
	}

	return nil
}

// GetFileInfo gets information about a file in GCS
func (g *GCSClient) GetFileInfo(ctx context.Context, objectName string) (*FileInfo, error) {
	obj := g.bucket.Object(objectName)
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get file attributes: %w", err)
	}

	return &FileInfo{
		Name:        attrs.Name,
		Size:        attrs.Size,
		ContentType: attrs.ContentType,
		Created:     attrs.Created,
		Updated:     attrs.Updated,
		URL:         fmt.Sprintf("https://storage.googleapis.com/%s/%s", g.bucketName, objectName),
		Bucket:      g.bucketName,
	}, nil
}

// ListFiles lists files in a GCS bucket with optional prefix
func (g *GCSClient) ListFiles(ctx context.Context, prefix string) ([]FileInfo, error) {
	var files []FileInfo

	query := &storage.Query{Prefix: prefix}
	it := g.bucket.Objects(ctx, query)

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate objects: %w", err)
		}

		files = append(files, FileInfo{
			Name:        attrs.Name,
			Size:        attrs.Size,
			ContentType: attrs.ContentType,
			Created:     attrs.Created,
			Updated:     attrs.Updated,
			URL:         fmt.Sprintf("https://storage.googleapis.com/%s/%s", g.bucketName, attrs.Name),
			Bucket:      g.bucketName,
		})
	}

	return files, nil
}

// DeleteFile deletes a file from Google Cloud Storage
func (g *GCSClient) DeleteFile(ctx context.Context, objectName string) error {
	obj := g.bucket.Object(objectName)
	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete file from GCS: %w", err)
	}

	return nil
}

// FileExists checks if a file exists in Google Cloud Storage
func (g *GCSClient) FileExists(ctx context.Context, objectName string) (bool, error) {
	obj := g.bucket.Object(objectName)
	_, err := obj.Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}

	return true, nil
}

// GenerateSignedURL generates a signed URL for temporary access to a file
func (g *GCSClient) GenerateSignedURL(ctx context.Context, objectName string, expiration time.Duration) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(expiration),
	}

	url, err := storage.SignedURL(g.bucketName, objectName, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return url, nil
}

// UploadJSON uploads JSON data to Google Cloud Storage
func (g *GCSClient) UploadJSON(ctx context.Context, objectName string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return g.UploadFile(ctx, objectName, jsonData, "application/json")
}

// DownloadJSON downloads and unmarshals JSON data from Google Cloud Storage
func (g *GCSClient) DownloadJSON(ctx context.Context, objectName string, target interface{}) error {
	data, err := g.DownloadFile(ctx, objectName)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

// Close closes the GCS client
func (g *GCSClient) Close() error {
	return g.client.Close()
}

// getContentTypeFromExtension determines content type from file extension
func getContentTypeFromExtension(ext string) string {
	switch strings.ToLower(ext) {
	case ".json":
		return "application/json"
	case ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".txt":
		return "text/plain"
	case ".csv":
		return "text/csv"
	case ".xml":
		return "application/xml"
	default:
		return "application/octet-stream"
	}
}

// Global GCS client instance
var gcsClient *GCSClient

// InitializeGCS initializes the global GCS client
func InitializeGCS() error {
	config := GCSConfig{
		ProjectID:       os.Getenv("GCS_PROJECT_ID"),
		BucketName:      os.Getenv("GCS_BUCKET_NAME"),
		CredentialsFile: os.Getenv("GCS_CREDENTIALS_FILE"),
	}

	if config.ProjectID == "" || config.BucketName == "" {
		return fmt.Errorf("GCS_PROJECT_ID and GCS_BUCKET_NAME environment variables are required")
	}

	client, err := NewGCSClient(config)
	if err != nil {
		return err
	}

	gcsClient = client
	return nil
}

// GetGCSClient returns the global GCS client
func GetGCSClient() *GCSClient {
	return gcsClient
}
