# Deployment Guide - Microservices Architecture

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                    API Gateway (kmqi5w0la0)                     │
│                        IAM Authentication                        │
└───────────────┬──────────────────────────┬─────────────────────┘
                │                          │
    ┌───────────▼──────────┐   ┌──────────▼──────────┐
    │  Document Service    │   │  Result Service     │
    │  /documents/*        │   │  /results/*         │
    └───────────┬──────────┘   └──────────┬──────────┘
                │                          │
    ┌───────────▼──────────┐   ┌──────────▼──────────┐
    │  Lambda Function     │   │  Lambda Function    │
    │  document-handler    │   │  result-handler     │
    └───────────┬──────────┘   └──────────┬──────────┘
                │                          │
    ┌───────────▼──────────┐   ┌──────────▼──────────┐
    │  S3 Bucket           │   │  Supabase Database  │
    │  security-...        │   │  (PostgreSQL)       │
    │  -document           │   │                     │
    └──────────────────────┘   └─────────────────────┘
                │
    ┌───────────▼──────────────────────────────────────┐
    │    Shared Deployment Bucket (S3)                 │
    │    security-questionnaire-deployment             │
    └──────────────────────────────────────────────────┘
```

## 📦 Resource Naming

### Document Service
- **Stack Name**: `security-questionnaire-document`
- **Lambda Function**: `security-questionnaire-document-handler`
- **S3 Bucket**: `security-questionnaire-document`
- **Deployment Prefix**: `document/`

### Result Service
- **Stack Name**: `security-questionnaire-result`
- **Lambda Function**: `security-questionnaire-result-handler`
- **Deployment Prefix**: `result/`

### Shared Resources
- **Deployment Bucket**: `security-questionnaire-deployment`
- **API Gateway ID**: `kmqi5w0la0`
- **Region**: `us-east-1`

## 🚀 Quick Start

### 1. Prerequisites Check

```bash
# Verify installations
go version             # Should be 1.21+
node --version         # Should be 18+
npm --version          # Should be 9+
serverless --version   # Should be 3.x
aws --version          # Should be 2.x
```

### 2. Set Environment Variables

```bash
export DATABASE_URL="postgresql://postgres:password@your-project.supabase.co:5432/postgres"
```

### 3. Install Dependencies

```bash
# Install Node.js dependencies (Serverless Framework & plugins)
make install
# or manually: npm install
```

### 4. Create Shared Resources

```bash
# Create deployment bucket (only needed once)
make create-bucket
```

### 5. Deploy All Services

```bash
# Deploy both microservices
make deploy-all
```

OR deploy individually:

```bash
# Deploy document service only
make deploy-document

# Deploy result service only
make deploy-result
```

## 📝 Step-by-Step Deployment

### Document Service

```bash
# Navigate to document service
cd services/document

# Deploy to development (auto-builds with serverless-go-plugin)
make deploy

# Deploy to production
make deploy-prod

# View service info
make info

# View logs
make logs
```

### Result Service

```bash
# Navigate to result service
cd services/result

# Deploy to development (auto-builds with serverless-go-plugin)
make deploy

# Deploy to production
make deploy-prod

# View service info
make info

# View logs
make logs
```

## 🔐 IAM Authentication

All API endpoints require AWS IAM authentication. Requests must be signed with AWS Signature Version 4.

### Testing with AWS CLI

```bash
# Install aws-api-gateway-cli-test
npm install -g aws-api-gateway-cli-test

# Test document endpoint
apig-test \
  --username 'YOUR_ACCESS_KEY' \
  --password 'YOUR_SECRET_KEY' \
  --invoke-url 'https://kmqi5w0la0.execute-api.us-east-1.amazonaws.com' \
  --api-gateway-name 'security-questionnaire' \
  --path-template '/documents' \
  --method 'GET' \
  --region 'us-east-1'
```

### Using Postman

1. Set request type (GET, POST, etc.)
2. Enter URL: `https://kmqi5w0la0.execute-api.us-east-1.amazonaws.com/documents`
3. Go to Authorization tab
4. Select Type: **AWS Signature**
5. Enter:
   - Access Key
   - Secret Key
   - AWS Region: `us-east-1`
   - Service Name: `execute-api`
6. Send request

## 🧪 Testing

### View Service Info

```bash
# Document service
cd services/document
make info

# Result service  
cd services/result
make info
```

### Test with Postman

Use the provided Postman collection in `postman/Document-Service.postman_collection.json`:

1. Import the collection
2. Set environment variables (API URL, AWS credentials)
3. Run requests with automatic IAM signing

See `postman/README.md` for detailed instructions.

## 📊 Monitoring

### View Lambda Logs

```bash
# Document service logs (from service directory)
cd services/document
make logs

# Or use AWS CLI directly
aws logs tail /aws/lambda/security-questionnaire-document-handler --follow

# Result service logs (from service directory)
cd services/result
make logs

# Or use AWS CLI directly
aws logs tail /aws/lambda/security-questionnaire-result-handler --follow
```

### CloudWatch Metrics

```bash
# View metrics in AWS Console
aws cloudwatch get-metric-statistics \
  --namespace AWS/Lambda \
  --metric-name Invocations \
  --dimensions Name=FunctionName,Value=security-questionnaire-document-handler \
  --start-time $(date -u -d '1 hour ago' +%Y-%m-%dT%H:%M:%S) \
  --end-time $(date -u +%Y-%m-%dT%H:%M:%S) \
  --period 300 \
  --statistics Sum
```

## 🗑️ Cleanup

### Delete Individual Services

```bash
cd services/document
make remove

cd services/result
make remove
```

### Delete All Services

```bash
# From project root
make delete-all
```

### Delete Deployment Bucket (Optional)

```bash
# Empty and delete deployment bucket
aws s3 rb s3://security-questionnaire-deployment --force
```

## 🔄 CI/CD Integration

### GitHub Actions Example

```yaml
name: Deploy Microservices

on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Setup Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.21
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      
      - name: Install Serverless Framework
        run: npm install
      
      - name: Configure AWS Credentials
        uses: aws-actions/configure-aws-credentials@v1
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-east-1
      
      - name: Deploy All Services
        env:
          DATABASE_URL: ${{ secrets.DATABASE_URL }}
        run: make deploy-all
```

## 📈 Scaling

### Increase Lambda Memory/Timeout

Edit `serverless.yml` in each service:

```yaml
functions:
  handler:
    timeout: 60          # Increase timeout (seconds)
    memorySize: 1024     # Increase memory (MB)
```

### Add Provisioned Concurrency

```yaml
functions:
  handler:
    provisionedConcurrency: 5
```

## 🆘 Troubleshooting

### Build Fails

```bash
# Clean and rebuild
cd services/document
make clean
go mod tidy
make deploy  # serverless-go-plugin will auto-build
```

### Deployment Fails

```bash
# Check AWS credentials
aws sts get-caller-identity

# Verify bucket exists
aws s3 ls s3://security-questionnaire-deployment

# Check CloudFormation stack
aws cloudformation describe-stacks \
  --stack-name security-questionnaire-document
```

### Import Path Errors

```bash
# Regenerate go.mod
cd project-root
rm go.sum
go mod tidy
```

## 💡 Best Practices

1. **Deploy to dev before prod** to catch issues
2. **Monitor CloudWatch logs** after deployment with `make logs`
3. **Use stage-specific configs** in serverless.yml
4. **Tag your resources** for better organization
5. **Set up CloudWatch alarms** for critical errors
6. **Implement retry logic** in your application code
7. **Use environment variables** for sensitive data
8. **Keep dependencies updated** with `npm audit` and `go mod tidy`

## 📚 Additional Resources

- [Serverless Framework Documentation](https://www.serverless.com/framework/docs/)
- [Serverless Go Plugin](https://github.com/mthenw/serverless-go-plugin)
- [API Gateway IAM Auth](https://docs.aws.amazon.com/apigateway/latest/developerguide/permissions.html)
- [Lambda Best Practices](https://docs.aws.amazon.com/lambda/latest/dg/best-practices.html)
- [Go Lambda Documentation](https://docs.aws.amazon.com/lambda/latest/dg/lambda-golang.html)
