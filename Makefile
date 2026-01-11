.PHONY: help install deploy-all deploy-infra deploy-document deploy-result delete-all delete-infra test deps clean

# Install dependencies
install:
	@echo "Installing Node.js dependencies..."
	@npm install
	@echo "Installing Go dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Dependencies installed"

# Create deployment bucket
create-bucket:
	@echo "Creating deployment bucket..."
	@aws s3 mb s3://security-questionnaire-deployment --region us-east-1 2>/dev/null || echo "Bucket already exists"
	@echo "✓ Bucket ready"

# Deploy infrastructure
deploy-infra: create-bucket
	@echo "Deploying shared infrastructure..."
	@npx serverless deploy --verbose
	@echo "✓ Infrastructure deployed"

# Deploy document service (serverless-go-plugin handles build & zip automatically)
deploy-document:
	@echo "Deploying Document Service..."
	@cd services/document && npx serverless deploy --verbose
	@echo "✓ Document Service deployed"

# Deploy result service (serverless-go-plugin handles build & zip automatically)
deploy-result:
	@echo "Deploying Result Service..."
	@cd services/result && npx serverless deploy --verbose
	@echo "✓ Result Service deployed"

# Deploy everything
deploy-all: install deploy-infra deploy-document deploy-result
	@echo ""
	@echo "================================================"
	@echo "✓ All services deployed successfully!"
	@echo "================================================"
	@npx serverless info --verbose 2>/dev/null | grep "HttpApiUrl" || echo ""

# Delete infrastructure
delete-infra:
	@echo "Deleting infrastructure..."
	@npx serverless remove --verbose
	@echo "✓ Infrastructure deleted"

# Delete all services
delete-all:
	@echo "Deleting all services..."
	@cd services/document && npx serverless remove --verbose || true
	@cd services/result && npx serverless remove --verbose || true
	@$(MAKE) delete-infra
	@echo "✓ All services deleted"

# Run tests
test:
	@echo "Running tests..."
	@go test ./... -v
	@echo "✓ Tests complete"

# Install/update Go dependencies
deps:
	@echo "Installing Go dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Go dependencies installed"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf .serverless
	@rm -rf services/document/.serverless
	@rm -rf services/result/.serverless
	@rm -rf node_modules
	@echo "✓ Clean complete"

# Show help
help:
	@echo "Security Questionnaire - Serverless Framework"
	@echo ""
	@echo "Setup:"
	@echo "  make install          - Install all dependencies"
	@echo ""
	@echo "Deploy:"
	@echo "  make deploy-all       - Deploy everything (infra + services)"
	@echo "  make deploy-infra     - Deploy shared infrastructure only"
	@echo "  make deploy-document  - Deploy document service only (auto-builds)"
	@echo "  make deploy-result    - Deploy result service only (auto-builds)"
	@echo ""
	@echo "Cleanup:"
	@echo "  make delete-all       - Delete all services"
	@echo "  make delete-infra     - Delete infrastructure only"
	@echo "  make clean            - Clean build artifacts"
	@echo ""
	@echo "Development:"
	@echo "  make test             - Run Go tests"
	@echo "  make deps             - Update Go dependencies"
	@echo ""
	@echo "Environment Variables Required:"
	@echo "  DATABASE_URL          - Supabase database connection string"
	@echo "  AWS credentials       - Configured via aws configure or env vars"
	@echo ""
	@echo "Note: serverless-go-plugin handles all Go building & zipping automatically!"
