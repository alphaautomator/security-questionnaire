package handlers

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// SuccessResponse creates a successful API Gateway response
func SuccessResponse(statusCode int, data interface{}) (events.APIGatewayProxyResponse, error) {
	body, _ := json.Marshal(data)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: string(body),
	}, nil
}

// ErrorResponse creates an error API Gateway response
func ErrorResponse(statusCode int, message string) (events.APIGatewayProxyResponse, error) {
	response := map[string]interface{}{
		"success": false,
		"message": message,
	}
	body, _ := json.Marshal(response)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: string(body),
	}, nil
}

// NotFoundResponse creates a 404 not found response
func NotFoundResponse() (events.APIGatewayProxyResponse, error) {
	return ErrorResponse(404, "Route not found")
}
