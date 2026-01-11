package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"security-questionnaire/config"
	"security-questionnaire/pkg/database"
	"security-questionnaire/services/document/models"

	"github.com/aws/aws-lambda-go/events"
)

// UpdateDocumentRequest represents the request body for updating a document
type UpdateDocumentRequest struct {
	Description *string `json:"description,omitempty"`
	Tags        *string `json:"tags,omitempty"`
}

// UpdateDocumentResponse represents the response for updating a document
type UpdateDocumentResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    *models.Document `json:"data,omitempty"`
}

// HandleUpdate handles updating a document's metadata
func HandleUpdate(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
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

	// Parse request body
	var req UpdateDocumentRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return ErrorResponse(400, "Invalid request body")
	}

	// Build updates map (only include fields that are provided)
	updates := make(map[string]interface{})
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Tags != nil {
		updates["tags"] = *req.Tags
	}

	if len(updates) == 0 {
		return ErrorResponse(400, "No fields to update")
	}

	// Initialize database service
	dbService, err := database.NewDatabaseService(cfg.DatabaseURL, &models.Document{})
	if err != nil {
		return ErrorResponse(500, fmt.Sprintf("Failed to initialize database service: %v", err))
	}
	defer dbService.Close()

	// Update document in database
	var doc models.Document
	if err := dbService.Update(&doc, documentID, updates); err != nil {
		return ErrorResponse(404, "Document not found or failed to update")
	}

	// Return success response
	response := UpdateDocumentResponse{
		Success: true,
		Message: "Document updated successfully",
		Data:    &doc,
	}

	return SuccessResponse(200, response)
}
