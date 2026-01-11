.PHONY: deploy deploy-prod remove test deps clean help

# Deploy to development stage
deploy:
	@echo "Deploying Result Service (dev)..."
	@cd ../.. && serverless deploy --config services/result/serverless.yml --stage dev
	@echo "Deployment complete!"

# Deploy to production stage
deploy-prod:
	@echo "Deploying Result Service (prod)..."
	@cd ../.. && serverless deploy --config services/result/serverless.yml --stage prod
	@echo "Production deployment complete!"

# Remove/delete stack
remove:
	@echo "Removing Result Service stack..."
	@cd ../.. && serverless remove --config services/result/serverless.yml --stage dev
	@echo "Stack removed!"

# Remove production stack
remove-prod:
	@echo "Removing Result Service production stack..."
	@cd ../.. && serverless remove --config services/result/serverless.yml --stage prod
	@echo "Production stack removed!"

# Run tests
test:
	@echo "Running tests..."
	@cd ../.. && go test ./services/result/... -v
	@echo "Tests complete!"

# Install Go dependencies
deps:
	@echo "Installing Go dependencies..."
	@cd ../.. && go mod download && go mod tidy
	@echo "Dependencies installed!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf .serverless .bin
	@echo "Clean complete!"

# Show service info
info:
	@echo "Getting Result Service info..."
	@cd ../.. && serverless info --config services/result/serverless.yml --stage dev

# View logs
logs:
	@echo "Fetching Result Service logs..."
	@cd ../.. && serverless logs --function handler --config services/result/serverless.yml --stage dev

# Show help
help:
	@echo "Result Service Commands:"
	@echo "  make deploy       - Deploy to development"
	@echo "  make deploy-prod  - Deploy to production"
	@echo "  make remove       - Remove development stack"
	@echo "  make remove-prod  - Remove production stack"
	@echo "  make test         - Run Go tests"
	@echo "  make deps         - Install dependencies"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make info         - Show service info"
	@echo "  make logs         - View function logs"
