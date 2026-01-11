package handlers

import (
	"context"
	"fmt"
	"time"

	"security-questionnaire/config"
	"security-questionnaire/pkg/database"
	"security-questionnaire/pkg/storage"
	"security-questionnaire/services/document/models"

	"github.com/aws/aws-lambda-go/events"
)

// ReadDocumentResponse represents the response for reading a document
type ReadDocumentResponse struct {
	Success      bool             `json:"success"`
	Message      string           `json:"message"`
	Data         *models.Document `json:"data,omitempty"`
	DownloadURL  string           `json:"download_url,omitempty"`
	URLExpiresIn string           `json:"url_expires_in,omitempty"`
}

// HandleRead handles reading a document by ID
func HandleRead(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return ErrorResponse(500, fmt.Sprintf("Configuration error: %v", err))
	}

	// Get document ID from path parameters
	documentID := request.PathParameters["id"]
	if documentID == "" {
		return ErrorResponse(400, "Document ID is required")
	}

	// Initialize database service
	dbService, err := database.NewDatabaseService(cfg.DatabaseURL, &models.Document{})
	if err != nil {
		return ErrorResponse(500, fmt.Sprintf("Failed to initialize database service: %v", err))
	}
	defer dbService.Close()

	// Get document from database
	var doc models.Document
	if err := dbService.GetByID(&doc, documentID); err != nil {
		return ErrorResponse(404, "Document not found")
	}

	// Initialize S3 service to generate pre-signed URL
	s3Service, err := storage.NewS3Service(cfg.S3Bucket, cfg.S3Region)
	if err != nil {
		return ErrorResponse(500, fmt.Sprintf("Failed to initialize S3 service: %v", err))
	}

	// Generate pre-signed URL (valid for 1 hour)
	downloadURL, err := s3Service.GetFileURL(doc.S3Key, 1*time.Hour)
	if err != nil {
		return ErrorResponse(500, fmt.Sprintf("Failed to generate download URL: %v", err))
	}

	// Return success response
	response := ReadDocumentResponse{
		Success:      true,
		Message:      "Document retrieved successfully",
		Data:         &doc,
		DownloadURL:  downloadURL,
		URLExpiresIn: "1 hour",
	}

	return SuccessResponse(200, response)
}
