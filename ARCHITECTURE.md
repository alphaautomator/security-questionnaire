# Architecture Overview

## 🏗️ Shared Library Structure

This project follows a **monorepo microservices architecture** with shared libraries to maximize code reuse and maintainability.

```
security-questionnaire/
│
├── pkg/                          # 🔥 Shared Libraries (Reusable)
│   ├── database/
│   │   └── database.go          # Generic GORM database service
│   ├── storage/
│   │   └── s3.go                # S3 storage service
│   └── models/
│       └── base.go              # Common base model
│
├── services/                     # Microservices
│   ├── document/
│   │   ├── cmd/api/main.go     # Entry point
│   │   ├── handlers/           # HTTP handlers
│   │   ├── models/             # Document-specific models
│   │   ├── serverless.yml      # Serverless Framework config
│   │   └── Makefile           # Service-specific commands
│   │
│   └── result/
│       ├── cmd/api/main.go     # Entry point
│       ├── handlers/           # HTTP handlers
│       ├── models/             # Result-specific models
│       ├── serverless.yml      # Serverless Framework config
│       └── Makefile           # Service-specific commands
│
├── config/                      # Shared configuration
├── go.mod                       # Single Go module for entire project
├── package.json                # Serverless Framework dependencies
├── serverless.yml              # Root serverless config (shared plugins)
└── Makefile                    # Build targets for all services
```

## 🎯 Key Design Decisions

### 1. Shared `pkg/` Directory

**Benefits:**
- ✅ **DRY Principle**: Write database/storage code once, use everywhere
- ✅ **Type Safety**: Generic functions with compile-time checks
- ✅ **Easy Maintenance**: Fix bugs in one place
- ✅ **Consistent Behavior**: All services use same database/storage logic

**Usage Example:**
```go
// In any service
import (
    "security-questionnaire/pkg/database"
    "security-questionnaire/pkg/storage"
)

// Initialize with any model
db, _ := database.NewDatabaseService(dbURL, &MyModel{})

// Generic operations
db.Create(model)
db.GetByID(model, id)
db.List(modelType, result, limit, offset)
```

### 2. Service-Specific Models

Each service has its own `models/` directory that:
- Extends the shared `BaseModel` from `pkg/models`
- Defines service-specific fields
- Keeps domain logic separate

**Example:**
```go
// services/document/models/document.go
type Document struct {
    models.BaseModel  // Inherits ID, CreatedAt, UpdatedAt, DeletedAt
    FileName    string
    FileSize    int64
    // ... document-specific fields
}

// services/result/models/result.go
type Result struct {
    models.BaseModel  // Inherits same base fields
    QuestionnaireID string
    Status          string
    // ... result-specific fields
}
```

### 3. Single Go Module

**Why?**
- Simplifies imports between `pkg/` and `services/`
- One `go.mod` to manage all dependencies
- Easier to build and test

**Build Process:**
- `serverless-go-plugin` automatically builds from service's `cmd/api` directory
- Plugin handles cross-compilation (Linux ARM64)
- Root `Makefile` orchestrates multi-service deployment
- All services have access to shared `pkg/` packages

### 4. Independent Deployment

Despite shared code:
- Each service has its own CloudFormation stack
- Services deploy independently
- Shared deployment bucket for efficiency

## 📦 Library Documentation

### `pkg/database` - Generic Database Service

**Features:**
- Generic CRUD operations for any GORM model
- Auto-migration support
- Pagination built-in
- Type-safe operations

**API:**
```go
// Initialize with models to auto-migrate
db := database.NewDatabaseService(url, &Model1{}, &Model2{})

// Create
db.Create(model)

// Read
db.GetByID(&model, id)

// List with pagination
total, _ := db.List(&ModelType{}, &results, limit, offset)

// Update
db.Update(&model, id, updates)

// Delete
db.Delete(&ModelType{}, id)

// Direct DB access for complex queries
db.GetDB().Where("custom = ?", val).Find(&results)
```

### `pkg/storage` - S3 Storage Service

**Features:**
- File upload with unique keys
- Pre-signed URL generation
- File download
- File deletion

**API:**
```go
s3 := storage.NewS3Service(bucket, region)

// Upload
key, url, _ := s3.UploadFile(storage.UploadFileData{
    FileName: "doc.pdf",
    FileContent: bytes,
    ContentType: "application/pdf",
})

// Get pre-signed URL
downloadURL, _ := s3.GetFileURL(key, 1*time.Hour)

// Download file
bytes, _ := s3.GetFile(key)

// Delete
s3.DeleteFile(key)
```

### `pkg/models` - Base Model

**Features:**
- Standard fields for all models
- GORM integration
- Soft delete support

**Definition:**
```go
type BaseModel struct {
    ID        string         `gorm:"primaryKey;type:uuid"`
    CreatedAt time.Time      `gorm:"column:created_at"`
    UpdatedAt time.Time      `gorm:"column:updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}
```

## 🔄 Request Flow

```
Client Request (IAM Signed)
    ↓
API Gateway (kmqi5w0la0)
    ↓
Lambda Function (Document/Result Handler)
    ↓
    ├─→ pkg/database → Supabase PostgreSQL
    └─→ pkg/storage  → AWS S3
```

## 🚀 Adding New Microservices

1. **Create service directory:**
```bash
mkdir services/new-service
```

2. **Import shared libraries:**
```go
import (
    "security-questionnaire/pkg/database"
    "security-questionnaire/pkg/storage"
    "security-questionnaire/pkg/models"
)
```

3. **Create service-specific model:**
```go
type NewModel struct {
    models.BaseModel
    // ... specific fields
}
```

4. **Use shared services:**
```go
db, _ := database.NewDatabaseService(dbURL, &NewModel{})
db.Create(&newModel)
```

5. **Add deploy target to root Makefile:**
```makefile
deploy-new-service:
    @echo "Deploying New Service..."
    @serverless deploy --config services/new-service/serverless.yml --stage $(STAGE)
```

6. **Create Serverless Framework config:**
```yaml
# services/new-service/serverless.yml
service: security-questionnaire-new-service

provider:
  name: aws
  runtime: provided.al2
  httpApi:
    id:
      Fn::ImportValue: ${self:custom.stackName}-HttpApiId

custom:
  stackName: security-questionnaire-infrastructure-${self:provider.stage}
  go:
    baseDir: cmd/api
    binDir: .bin
    monorepo: true

plugins:
  - serverless-go-plugin

functions:
  handler:
    name: security-questionnaire-new-service-handler
    handler: bootstrap
    events:
      - httpApi:
          path: /new-service
          method: GET
```

7. **Deploy!**
```bash
cd services/new-service
make deploy
```

## 📊 Benefits Summary

| Feature | Traditional Approach | Our Approach |
|---------|---------------------|--------------|
| Code Reuse | Duplicate in each service | Shared `pkg/` libraries |
| Maintenance | Fix bugs in N places | Fix once in `pkg/` |
| Type Safety | Manual type checking | Generic with compile-time checks |
| Build Complexity | Simple per-service | Single root build |
| Testing | Test each service separately | Test shared libs + services |
| Onboarding | Learn each service | Learn `pkg/` once |

## 🔐 Security

- IAM authentication on all endpoints
- Shared libraries follow AWS best practices
- Minimal permissions per service
- Secure S3 operations with pre-signed URLs
- Database connections use SSL

## 📈 Scalability

- Add new services without duplicating code
- Scale services independently
- Shared libraries tested once, work everywhere
- Easy to add new features to `pkg/`

## 🧪 Testing Strategy

```bash
# Test shared libraries
go test ./pkg/... -v

# Test specific service
go test ./services/document/... -v

# Test everything
go test ./... -v
```

## 💡 Best Practices

1. **Keep `pkg/` generic** - No service-specific logic
2. **Service models extend BaseModel** - Consistency
3. **Use generic database operations** - Type safety
4. **Build from root** - Access to all packages
5. **One dependency version** - Single `go.mod`
6. **Document shared APIs** - Help other developers

## 🎓 Learning Path

For new developers:

1. Start with `pkg/models/base.go` - Understand base model
2. Read `pkg/database/database.go` - Learn generic CRUD
3. Read `pkg/storage/s3.go` - Learn S3 operations
4. Explore `services/document` - See usage example
5. Create new service using patterns

This architecture scales from 2 to 200 microservices! 🚀
