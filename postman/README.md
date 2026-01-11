# Postman Collection for Security Questionnaire

This directory contains Postman collections for testing the Security Questionnaire microservices.

## 📦 Files

- **`Document-Service.postman_collection.json`** - Complete API collection for Document Service
- **`Security-Questionnaire.postman_environment.json`** - Environment variables

## 🚀 Quick Start

### 1. Import into Postman

**Import Collection:**
1. Open Postman
2. Click **Import** button (top left)
3. Select `Document-Service.postman_collection.json`
4. Click **Import**

**Import Environment:**
1. Click **Import** button
2. Select `Security-Questionnaire.postman_environment.json`
3. Click **Import**
4. Select the environment from the dropdown (top right)

### 2. Configure AWS Credentials

**Option A: Edit Environment Variables**
1. Click the **eye icon** (top right) → **Edit**
2. Update the following variables:
   - `aws_access_key`: Your AWS Access Key ID
   - `aws_secret_key`: Your AWS Secret Access Key
3. Click **Save**

**Option B: Use Collection Authorization**
1. Click on the collection name
2. Go to **Authorization** tab
3. Type: **AWS Signature**
4. Fill in:
   - AccessKey: Your AWS Access Key ID
   - SecretKey: Your AWS Secret Access Key
   - AWS Region: `us-east-1`
   - Service Name: `execute-api`

### 3. Test the API

**Test Sequence:**
1. **Create Document** → Creates a document and auto-saves the ID
2. **List Documents** → View all documents
3. **Get Document** → Uses saved document ID
4. **Update Document** → Updates metadata
5. **Delete Document** → Removes document

## 📝 Available Endpoints

### Document Service

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/documents` | Upload new document |
| GET | `/documents` | List all documents (paginated) |
| GET | `/documents/{id}` | Get specific document |
| PUT | `/documents/{id}` | Update document metadata |
| DELETE | `/documents/{id}` | Delete document |

## 🔑 AWS IAM Authentication

All endpoints use **AWS Signature Version 4** authentication.

### Required IAM Permissions

Your AWS user/role needs:
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "execute-api:Invoke",
      "Resource": "arn:aws:execute-api:us-east-1:*:kmqi5w0la0/*"
    }
  ]
}
```

## 📄 Request Examples

### Create Document

```json
{
  "file_name": "security-report.pdf",
  "file_content": "BASE64_ENCODED_CONTENT_HERE",
  "content_type": "application/pdf",
  "description": "Q4 2023 Security Assessment",
  "tags": "security, assessment, q4"
}
```

**How to get base64 content:**

**macOS/Linux:**
```bash
base64 -i myfile.pdf | tr -d '\n'
```

**Windows (PowerShell):**
```powershell
[Convert]::ToBase64String([IO.File]::ReadAllBytes("myfile.pdf"))
```

**JavaScript:**
```javascript
// In Pre-request Script tab
const fs = require('fs');
const fileContent = fs.readFileSync('/path/to/file.pdf');
const base64 = fileContent.toString('base64');
pm.environment.set('base64_file', base64);
```

### Update Document

```json
{
  "description": "Updated security assessment report",
  "tags": "security, updated, final"
}
```

## 🎯 Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `base_url` | API Gateway URL | `https://kmqi5w0la0.execute-api.us-east-1.amazonaws.com` |
| `aws_access_key` | AWS Access Key ID | `AKIAIOSFODNN7EXAMPLE` |
| `aws_secret_key` | AWS Secret Access Key | `wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY` |
| `aws_region` | AWS Region | `us-east-1` |
| `document_id` | Auto-set after creating a document | `123e4567-e89b-12d3-a456-426614174000` |
| `base64_sample_file` | Sample base64 content for testing | Auto-generated |

## 🧪 Testing Features

### Auto-save Document ID

After creating a document, the collection automatically saves the `document_id` to the environment. This allows you to:
1. Create a document
2. Immediately test Get/Update/Delete without manual ID copying

### Pre-request Scripts

- Generates sample base64 content if not provided
- Sets up required headers

### Test Scripts

- Auto-saves document ID from Create response
- Logs response for debugging
- Validates status codes

## 🐛 Troubleshooting

### Error: 403 Forbidden

**Cause:** Invalid AWS credentials or insufficient permissions

**Solution:**
1. Verify AWS credentials are correct
2. Check IAM permissions for `execute-api:Invoke`
3. Ensure credentials belong to the correct AWS account

### Error: 400 Bad Request - "Invalid base64"

**Cause:** File content is not properly base64 encoded

**Solution:**
1. Ensure file content is base64 encoded
2. Remove any line breaks or whitespace
3. Use `-w 0` flag with base64 command: `base64 -w 0 file.pdf`

### Error: 404 Not Found

**Cause:** Document ID doesn't exist or was deleted

**Solution:**
1. Run "List Documents" to see available documents
2. Copy a valid document ID
3. Update the `document_id` variable

### Error: 500 Internal Server Error

**Possible causes:**
1. Database connection issue (check Supabase)
2. S3 bucket doesn't exist
3. Lambda function error

**Solution:**
1. Check CloudWatch logs:
   ```bash
   aws logs tail /aws/lambda/security-questionnaire-document-handler --follow
   ```
2. Verify DATABASE_URL environment variable
3. Ensure S3 bucket exists

## 📊 Response Examples

### Success Response (Create Document)

```json
{
  "success": true,
  "message": "Document created successfully",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "file_name": "security-report.pdf",
    "file_size": 204800,
    "content_type": "application/pdf",
    "s3_bucket": "security-questionnaire-document",
    "s3_key": "documents/abc123.pdf",
    "s3_url": "https://security-questionnaire-document.s3.amazonaws.com/documents/abc123.pdf",
    "description": "Q4 2023 Security Assessment",
    "tags": "security, assessment, q4",
    "created_at": "2026-01-11T10:30:00Z",
    "updated_at": "2026-01-11T10:30:00Z"
  }
}
```

### Success Response (Get Document)

```json
{
  "success": true,
  "message": "Document retrieved successfully",
  "data": { ... },
  "download_url": "https://security-questionnaire-document.s3.amazonaws.com/documents/abc123.pdf?X-Amz-...",
  "url_expires_in": "1 hour"
}
```

### Error Response

```json
{
  "success": false,
  "message": "Document not found"
}
```

## 🔄 Workflow Example

### Complete CRUD Workflow

1. **Create** a document
   - Request: POST `/documents` with file data
   - Response: Document created with ID
   - Auto-saved: `document_id` variable

2. **List** all documents
   - Request: GET `/documents?limit=10&offset=0`
   - Response: Array of documents

3. **Get** specific document
   - Request: GET `/documents/{{document_id}}`
   - Response: Document details + download URL

4. **Update** document metadata
   - Request: PUT `/documents/{{document_id}}`
   - Response: Updated document

5. **Delete** document
   - Request: DELETE `/documents/{{document_id}}`
   - Response: Success confirmation

## 💡 Tips

1. **Use Environment Variables**: Switch between dev/prod easily
2. **Save Responses**: Use Postman's save response feature
3. **Create Test Scripts**: Validate responses automatically
4. **Use Pre-request Scripts**: Generate dynamic data
5. **Organize Folders**: Group related requests
6. **Share Collection**: Export and share with team

## 🔐 Security Best Practices

1. **Never commit AWS credentials** to version control
2. **Use IAM roles** in production instead of access keys
3. **Rotate credentials** regularly
4. **Use temporary credentials** (STS) when possible
5. **Set minimal permissions** (principle of least privilege)

## 📚 Additional Resources

- [Postman AWS Signature](https://learning.postman.com/docs/sending-requests/authorization/#aws-signature)
- [AWS Signature V4](https://docs.aws.amazon.com/general/latest/gr/signature-version-4.html)
- [API Gateway IAM Auth](https://docs.aws.amazon.com/apigateway/latest/developerguide/permissions.html)

## 🆘 Support

If you encounter issues:
1. Check CloudWatch logs for Lambda errors
2. Verify API Gateway is deployed
3. Confirm IAM permissions
4. Test with AWS CLI first
5. Review Postman console for detailed errors
