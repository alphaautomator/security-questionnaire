package handlers

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// SuccessResponse creates a successful API Gateway V2 response
func SuccessResponse(statusCode int, data interface{}) (events.APIGatewayV2HTTPResponse, error) {
	body, _ := json.Marshal(data)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: string(body),
	}, nil
}

// ErrorResponse creates an error API Gateway V2 response
func ErrorResponse(statusCode int, message string) (events.APIGatewayV2HTTPResponse, error) {
	response := map[string]interface{}{
		"success": false,
		"message": message,
	}
	body, _ := json.Marshal(response)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: string(body),
	}, nil
}

// NotFoundResponse creates a 404 not found response
func NotFoundResponse() (events.APIGatewayV2HTTPResponse, error) {
	return ErrorResponse(404, "Route not found")
}
