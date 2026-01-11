package handlers

import (
	"context"
	"fmt"
	"strconv"

	"security-questionnaire/config"
	"security-questionnaire/pkg/database"
	"security-questionnaire/services/document/models"

	"github.com/aws/aws-lambda-go/events"
)

// ListDocumentsResponse represents the response for listing documents
type ListDocumentsResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    []models.Document `json:"data"`
	Total   int64             `json:"total"`
	Limit   int               `json:"limit"`
	Offset  int               `json:"offset"`
}

// HandleList handles listing all documents with pagination
func HandleList(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return ErrorResponse(500, fmt.Sprintf("Configuration error: %v", err))
	}

	// Parse pagination parameters
	limit := 10 // default
	offset := 0 // default

	if limitStr := request.QueryStringParameters["limit"]; limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr := request.QueryStringParameters["offset"]; offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Initialize database service
	dbService, err := database.NewDatabaseService(cfg.DatabaseURL, &models.Document{})
	if err != nil {
		return ErrorResponse(500, fmt.Sprintf("Failed to initialize database service: %v", err))
	}
	defer dbService.Close()

	// Get documents from database
	var documents []models.Document
	total, err := dbService.List(&models.Document{}, &documents, limit, offset)
	if err != nil {
		return ErrorResponse(500, fmt.Sprintf("Failed to list documents: %v", err))
	}

	// Return success response
	response := ListDocumentsResponse{
		Success: true,
		Message: "Documents retrieved successfully",
		Data:    documents,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	}

	return SuccessResponse(200, response)
}
