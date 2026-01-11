# Security Questionnaire - Microservices Backend

A serverless microservices backend built with **Go**, **AWS Lambda**, **API Gateway**, **S3**, and **Supabase** (PostgreSQL).  
Deployed using **Serverless Framework**.

## 🏗️ Architecture

```
Infrastructure Stack (API Gateway)
├── Document Service
│   ├── Lambda (Go)
│   ├── S3 Bucket
│   └── Routes: /documents, /documents/{id}
└── Result Service
    ├── Lambda (Go)
    └── Routes: /results, /results/{id}
```

## 📂 Project Structure

```
security-questionnaire/
├── serverless.yml           # Infrastructure (API Gateway)
├── package.json             # Node.js dependencies (Serverless Framework)
├── Makefile                 # Build & deployment commands
│
├── services/
│   ├── document/
│   │   ├── serverless.yml   # Document service config
│   │   ├── cmd/api/main.go  # Lambda entry point
│   │   ├── handlers/        # Request handlers
│   │   └── models/          # Domain models
│   │
│   └── result/
│       ├── serverless.yml   # Result service config
│       ├── cmd/api/main.go  # Lambda entry point
│       ├── handlers/
│       └── models/
│
└── pkg/                     # Shared libraries
    ├── database/            # Generic GORM database service
    ├── storage/             # S3 storage utilities
    └── models/              # Shared base models
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- AWS CLI configured
- Supabase account & database URL

### 1. Install Dependencies

```bash
make install
```

This installs both Node.js (Serverless Framework) and Go dependencies.

### 2. Set Environment Variables

```bash
export DATABASE_URL="postgresql://user:pass@host:port/db?sslmode=require"
```

### 3. Deploy Everything

```bash
make deploy-all
```

This will:
1. Install dependencies
2. Deploy API Gateway (infrastructure)
3. Build & deploy Document Service
4. Build & deploy Result Service
5. Output the API endpoint

## 📦 Deployment Commands

### Full Deployment

```bash
make deploy-all              # Deploy infrastructure + all services
```

### Individual Components

```bash
make install                 # Install Node.js & Go dependencies
make deploy-infra            # Deploy API Gateway only
make deploy-document         # Deploy document service only
make deploy-result           # Deploy result service only
```

### Build Only

```bash
make build-document          # Build document handler
make build-result            # Build result handler
```

### Cleanup

```bash
make delete-all              # Delete all stacks
make delete-infra            # Delete infrastructure only
make clean                   # Clean build artifacts
```

### Development

```bash
make test                    # Run all tests
make deps                    # Update Go dependencies
make help                    # Show all commands
```

## 🔌 API Endpoints

After deployment, you'll get an API endpoint like:
```
https://xxxxx.execute-api.us-east-1.amazonaws.com/dev
```

### Document Service

| Method | Path | Description |
|--------|------|-------------|
| POST | `/documents` | Create a new document |
| GET | `/documents` | List all documents (paginated) |
| GET | `/documents/{id}` | Get document by ID |
| PUT | `/documents/{id}` | Update document metadata |
| DELETE | `/documents/{id}` | Delete document |

### Result Service

| Method | Path | Description |
|--------|------|-------------|
| POST | `/results` | Create a new result |
| GET | `/results` | List all results (paginated) |
| GET | `/results/{id}` | Get result by ID |
| PUT | `/results/{id}` | Update result |
| DELETE | `/results/{id}` | Delete result |

## 🔐 Authentication

All endpoints use **AWS IAM Authentication**.

### Using AWS CLI

```bash
aws apigatewayv2 invoke \
  --api-id xxxxx \
  --stage dev \
  --request-items '{
    "httpMethod": "GET",
    "path": "/documents"
  }'
```

### Using Postman

See `postman/` directory for:
- Collection with all API endpoints
- Environment variables template
- AWS IAM authentication setup

## 🛠️ Technology Stack

- **Language**: Go 1.21
- **Runtime**: AWS Lambda (provided.al2, ARM64)
- **ORM**: GORM
- **Database**: Supabase (PostgreSQL)
- **Storage**: AWS S3
- **API Gateway**: AWS HTTP API (API Gateway V2)
- **IaC**: Serverless Framework
- **Authentication**: AWS IAM

## 📊 Database Schema

### Documents Table

```sql
CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_name VARCHAR NOT NULL,
    file_size BIGINT NOT NULL,
    content_type VARCHAR NOT NULL,
    s3_bucket VARCHAR NOT NULL,
    s3_key VARCHAR NOT NULL,
    description TEXT,
    tags TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

## 🔄 How It Works

### Serverless Framework Benefits

✅ **Simpler Config**: YAML-based, cleaner than SAM
✅ **Automatic Packaging**: Handles ZIP creation
✅ **CloudFormation Exports**: Easy cross-stack references
✅ **Plugin Ecosystem**: Rich plugins available
✅ **Better DX**: Faster deployments, better error messages

### Deployment Flow

1. **Infrastructure Stack** creates shared API Gateway
2. **Each service**:
   - Builds Go binary (ARM64)
   - Packages as ZIP
   - Deploys to Lambda
   - Automatically creates routes on shared API Gateway
   - Uses CloudFormation exports for API Gateway ID

### Shared Libraries

- `pkg/database`: Generic GORM database service (type-safe with generics)
- `pkg/storage`: S3 operations (upload, download, delete, presigned URLs)
- `pkg/models`: Common base model with ID, timestamps, soft delete

## 🧪 Local Development

### Run Tests

```bash
make test
```

### Local Invocation

```bash
# Install serverless-offline
npm install --save-dev serverless-offline

# Add to serverless.yml plugins
# Run locally
npx serverless offline
```

## 📝 Adding a New Microservice

1. Create service directory:
```bash
mkdir -p services/newservice/{cmd/api,handlers,models}
```

2. Copy serverless config:
```bash
cp services/document/serverless.yml services/newservice/
```

3. Update `services/newservice/serverless.yml`:
   - Change `service` name
   - Update function name
   - Define your routes

4. Implement Go handlers in `services/newservice/handlers/`

5. Add build/deploy targets to root `Makefile`

6. Deploy:
```bash
make build-newservice
cd services/newservice && npx serverless deploy
```

## 🎯 Why Serverless Framework?

| Feature | SAM | Serverless Framework |
|---------|-----|---------------------|
| Config Syntax | Verbose | Clean & Simple |
| CloudFormation Exports | Manual | Automatic |
| Plugin Ecosystem | Limited | Rich |
| Deployment Speed | Slower | Faster |
| Error Messages | Cryptic | Clear |
| Local Testing | sam local | serverless offline |

## 🐛 Troubleshooting

### Build Fails

Make sure you're building for Linux ARM64:
```bash
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build
```

### Database Connection Issues

Check your `DATABASE_URL` format:
```bash
postgresql://user:password@host:port/database?sslmode=require
```

### Deployment Fails

1. Check AWS credentials: `aws sts get-caller-identity`
2. Verify bucket exists: `aws s3 ls s3://security-questionnaire-deployment`
3. Check CloudFormation stack: `npx serverless info`

### API Gateway Not Found

Deploy infrastructure first:
```bash
make deploy-infra
```

## 📄 License

MIT

## 🤝 Contributing

1. Create a feature branch
2. Make your changes
3. Run tests: `make test`
4. Deploy to dev: `make deploy-all`
5. Submit a pull request

---

Built with ❤️ using Go + Serverless Framework
