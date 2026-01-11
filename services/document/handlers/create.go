package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"security-questionnaire/config"
	"security-questionnaire/pkg/database"
	"security-questionnaire/pkg/storage"
	"security-questionnaire/services/document/models"

	"github.com/aws/aws-lambda-go/events"
)

// CreateDocumentRequest represents the request body for creating a document
type CreateDocumentRequest struct {
	FileName    string `json:"file_name"`
	FileContent string `json:"file_content"` // base64 encoded
	ContentType string `json:"content_type"`
	Description string `json:"description,omitempty"`
	Tags        string `json:"tags,omitempty"`
}

// CreateDocumentResponse represents the response for creating a document
type CreateDocumentResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    *models.Document `json:"data,omitempty"`
}

// HandleCreate handles the creation of a new document
func HandleCreate(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return ErrorResponse(500, fmt.Sprintf("Configuration error: %v", err))
	}

	// Parse request body
	var req CreateDocumentRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return ErrorResponse(400, "Invalid request body")
	}

	// Validate required fields
	if req.FileName == "" || req.FileContent == "" || req.ContentType == "" {
		return ErrorResponse(400, "file_name, file_content, and content_type are required")
	}

	// Decode base64 file content
	fileBytes, err := base64.StdEncoding.DecodeString(req.FileContent)
	if err != nil {
		return ErrorResponse(400, "Invalid base64 encoded file content")
	}

	// Initialize S3 service
	s3Service, err := storage.NewS3Service(cfg.S3Bucket, cfg.S3Region)
	if err != nil {
		return ErrorResponse(500, fmt.Sprintf("Failed to initialize S3 service: %v", err))
	}

	// Upload file to S3
	s3Key, s3URL, err := s3Service.UploadFile(storage.UploadFileData{
		FileName:    req.FileName,
		FileContent: fileBytes,
		ContentType: req.ContentType,
	})
	if err != nil {
		return ErrorResponse(500, fmt.Sprintf("Failed to upload file: %v", err))
	}

	// Initialize database service
	dbService, err := database.NewDatabaseService(cfg.DatabaseURL, &models.Document{})
	if err != nil {
		// Cleanup: delete uploaded file from S3
		_ = s3Service.DeleteFile(s3Key)
		return ErrorResponse(500, fmt.Sprintf("Failed to initialize database service: %v", err))
	}
	defer dbService.Close()

	// Create document record in database
	doc := &models.Document{
		FileName:    req.FileName,
		FileSize:    int64(len(fileBytes)),
		ContentType: req.ContentType,
		S3Bucket:    cfg.S3Bucket,
		S3Key:       s3Key,
		S3URL:       s3URL,
		Description: req.Description,
		Tags:        req.Tags,
	}

	if err := dbService.Create(doc); err != nil {
		// Cleanup: delete uploaded file from S3
		_ = s3Service.DeleteFile(s3Key)
		return ErrorResponse(500, fmt.Sprintf("Failed to create document record: %v", err))
	}

	// Return success response
	response := CreateDocumentResponse{
		Success: true,
		Message: "Document created successfully",
		Data:    doc,
	}

	return SuccessResponse(201, response)
}
