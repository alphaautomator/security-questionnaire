package main

import (
	"context"
	"security-questionnaire/services/result/handlers"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Router handles all API requests and routes them to appropriate handlers
func Router(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// Route based on HTTP method and path (HTTP API V2 format)
	method := request.RequestContext.HTTP.Method
	path := request.RawPath
	if path == "" {
		path = request.RequestContext.HTTP.Path
	}

	// Handle different routes
	switch {
	case method == "POST" && path == "/results":
		return handlers.HandleCreate(ctx, request)

	case method == "GET" && path == "/results":
		return handlers.HandleList(ctx, request)

	case method == "GET" && request.PathParameters["id"] != "":
		return handlers.HandleRead(ctx, request)

	case method == "PUT" && request.PathParameters["id"] != "":
		return handlers.HandleUpdate(ctx, request)

	case method == "DELETE" && request.PathParameters["id"] != "":
		return handlers.HandleDelete(ctx, request)

	default:
		return handlers.NotFoundResponse()
	}
}

func main() {
	lambda.Start(Router)
}
