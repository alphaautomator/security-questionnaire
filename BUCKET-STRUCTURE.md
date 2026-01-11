# Deployment Bucket Structure

## 📦 Shared Deployment Bucket

**Bucket Name**: `security-questionnaire-deployment`

All microservices use this **single shared S3 bucket** for storing deployment artifacts (CloudFormation templates, Lambda code packages, etc.).

## 🗂️ Bucket Organization

```
s3://security-questionnaire-deployment/
│
├── document/                          # Document service artifacts (dev)
│   ├── template.yaml                  # Packaged CloudFormation template
│   └── [build-id]/                    # Lambda code packages
│       └── bootstrap.zip
│
├── document-prod/                     # Document service artifacts (prod)
│   ├── template.yaml
│   └── [build-id]/
│       └── bootstrap.zip
│
├── result/                            # Result service artifacts (dev)
│   ├── template.yaml
│   └── [build-id]/
│       └── bootstrap.zip
│
└── result-prod/                       # Result service artifacts (prod)
    ├── template.yaml
    └── [build-id]/
        └── bootstrap.zip
```

## 🎯 Why One Bucket?

### Benefits:
- ✅ **Cost Efficient**: Single bucket to manage
- ✅ **Organized**: Clear prefixes per service
- ✅ **Simple Permissions**: One bucket policy
- ✅ **Easy Cleanup**: Delete one bucket to clean all
- ✅ **Separation**: Each service has its own prefix

### Prefix Strategy:
- **Development**: `{service-name}/` (e.g., `document/`, `result/`)
- **Production**: `{service-name}-prod/` (e.g., `document-prod/`, `result-prod/`)

## 🚀 Usage

### Create Bucket (One Time)
```bash
make create-bucket
# or
aws s3 mb s3://security-questionnaire-deployment --region us-east-1
```

### Deploy Services
```bash
# Document service uploads to: s3://security-questionnaire-deployment/document/
make deploy-document

# Result service uploads to: s3://security-questionnaire-deployment/result/
make deploy-result
```

### Production Deployment
```bash
# Document service uploads to: s3://security-questionnaire-deployment/document-prod/
cd services/document && make deploy-prod

# Result service uploads to: s3://security-questionnaire-deployment/result-prod/
cd services/result && make deploy-prod
```

## 🔍 View Bucket Contents

```bash
# List all prefixes
aws s3 ls s3://security-questionnaire-deployment/

# List document service artifacts
aws s3 ls s3://security-questionnaire-deployment/document/ --recursive

# List result service artifacts
aws s3 ls s3://security-questionnaire-deployment/result/ --recursive
```

## 🗑️ Cleanup

### Clean Service Artifacts
```bash
# Remove document service artifacts
aws s3 rm s3://security-questionnaire-deployment/document/ --recursive

# Remove result service artifacts
aws s3 rm s3://security-questionnaire-deployment/result/ --recursive
```

### Delete Entire Bucket
```bash
# Remove all artifacts and delete bucket
aws s3 rb s3://security-questionnaire-deployment --force
```

## 📊 Bucket Policy

The bucket is private and only accessible by:
- AWS SAM CLI during deployment
- CloudFormation for retrieving templates
- Lambda for fetching code packages

## 🔐 Security

- **Private Access**: No public access
- **IAM Controlled**: Only authenticated AWS principals
- **Regional**: Created in `us-east-1`
- **Versioning**: Optional (can be enabled for rollback)

## 💰 Cost Optimization

### Current Setup:
- **Storage**: Pay only for actual artifacts
- **Requests**: Minimal (only during deployment)
- **Transfer**: Free within same region

### Recommendations:
1. **Lifecycle Policy**: Auto-delete old build artifacts after 30 days
2. **Intelligent Tiering**: Move infrequent artifacts to cheaper storage
3. **Versioning**: Enable only if rollback is critical

### Add Lifecycle Policy:
```bash
aws s3api put-bucket-lifecycle-configuration \
  --bucket security-questionnaire-deployment \
  --lifecycle-configuration file://lifecycle-policy.json
```

**lifecycle-policy.json**:
```json
{
  "Rules": [
    {
      "Id": "DeleteOldArtifacts",
      "Status": "Enabled",
      "Prefix": "",
      "Expiration": {
        "Days": 30
      }
    }
  ]
}
```

## 📈 Adding New Services

When adding a new microservice:

1. **Update samconfig.toml**:
```toml
[default.deploy.parameters]
s3_bucket = "security-questionnaire-deployment"
s3_prefix = "new-service"  # ← Your service name
```

2. **Production config**:
```toml
[prod.deploy.parameters]
s3_bucket = "security-questionnaire-deployment"
s3_prefix = "new-service-prod"  # ← Your service name with -prod
```

3. **Deploy**: Artifacts automatically go to the right prefix!

## 🎓 Best Practices

1. ✅ **Use consistent prefixes**: `{service-name}/` for dev, `{service-name}-prod/` for prod
2. ✅ **One bucket per environment** (optional): Or use prefixes as we do
3. ✅ **Enable logging**: Track deployment activities
4. ✅ **Tag the bucket**: For cost allocation
5. ✅ **Monitor size**: Set up CloudWatch alarms for unusual growth

## 📝 Bucket Tags (Recommended)

```bash
aws s3api put-bucket-tagging \
  --bucket security-questionnaire-deployment \
  --tagging 'TagSet=[
    {Key=Project,Value=security-questionnaire},
    {Key=Purpose,Value=deployment-artifacts},
    {Key=ManagedBy,Value=sam-cli}
  ]'
```

This ensures all your deployment artifacts are organized, efficient, and easy to manage! 🚀
