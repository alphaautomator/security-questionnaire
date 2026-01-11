.PHONY: deploy deploy-prod remove test deps clean info logs help

# Deploy to development stage
deploy:
	@echo "Deploying Document Service (dev)..."
	@cd ../.. && serverless deploy --config services/document/serverless.yml --stage dev
	@echo "Deployment complete!"

# Deploy to production stage
deploy-prod:
	@echo "Deploying Document Service (prod)..."
	@cd ../.. && serverless deploy --config services/document/serverless.yml --stage prod
	@echo "Production deployment complete!"

# Remove/delete stack
remove:
	@echo "Removing Document Service stack..."
	@cd ../.. && serverless remove --config services/document/serverless.yml --stage dev
	@echo "Stack removed!"

# Remove production stack
remove-prod:
	@echo "Removing Document Service production stack..."
	@cd ../.. && serverless remove --config services/document/serverless.yml --stage prod
	@echo "Production stack removed!"

# Run tests
test:
	@echo "Running tests..."
	@cd ../.. && go test ./services/document/... -v
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
	@echo "Getting Document Service info..."
	@cd ../.. && serverless info --config services/document/serverless.yml --stage dev

# View logs
logs:
	@echo "Fetching Document Service logs..."
	@cd ../.. && serverless logs --function handler --config services/document/serverless.yml --stage dev

# Show help
help:
	@echo "Document Service Commands:"
	@echo "  make deploy       - Deploy to development"
	@echo "  make deploy-prod  - Deploy to production"
	@echo "  make remove       - Remove development stack"
	@echo "  make remove-prod  - Remove production stack"
	@echo "  make test         - Run Go tests"
	@echo "  make deps         - Install dependencies"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make info         - Show service info"
	@echo "  make logs         - View function logs"
