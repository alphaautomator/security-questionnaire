package handlers

import (
	"context"
	"fmt"

	"security-questionnaire/config"
	"security-questionnaire/pkg/database"
	"security-questionnaire/pkg/storage"
	"security-questionnaire/services/document/models"

	"github.com/aws/aws-lambda-go/events"
)

// DeleteDocumentResponse represents the response for deleting a document
type DeleteDocumentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// HandleDelete handles deleting a document by ID
func HandleDelete(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
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

	// Get document details before deletion (to get S3 key)
	var doc models.Document
	if err := dbService.GetByID(&doc, documentID); err != nil {
		return ErrorResponse(404, "Document not found")
	}

	// Initialize S3 service
	s3Service, err := storage.NewS3Service(cfg.S3Bucket, cfg.S3Region)
	if err != nil {
		return ErrorResponse(500, fmt.Sprintf("Failed to initialize S3 service: %v", err))
	}

	// Delete file from S3
	if err := s3Service.DeleteFile(doc.S3Key); err != nil {
		return ErrorResponse(500, fmt.Sprintf("Failed to delete file from S3: %v", err))
	}

	// Delete document from database
	if err := dbService.Delete(&models.Document{}, documentID); err != nil {
		return ErrorResponse(500, fmt.Sprintf("Failed to delete document: %v", err))
	}

	// Return success response
	response := DeleteDocumentResponse{
		Success: true,
		Message: "Document deleted successfully",
	}

	return SuccessResponse(200, response)
}
