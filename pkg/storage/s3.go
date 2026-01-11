package storage

import (
	"bytes"
	"fmt"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

// S3Service handles all S3 operations
type S3Service struct {
	client   *s3.S3
	uploader *s3manager.Uploader
	bucket   string
	region   string
}

// NewS3Service creates a new S3 service
func NewS3Service(bucket, region string) (*S3Service, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	return &S3Service{
		client:   s3.New(sess),
		uploader: s3manager.NewUploader(sess),
		bucket:   bucket,
		region:   region,
	}, nil
}

// UploadFileData uploads file data to S3
type UploadFileData struct {
	FileName    string
	FileContent []byte
	ContentType string
}

// UploadFile uploads a file to S3 and returns the S3 key and URL
func (s *S3Service) UploadFile(data UploadFileData) (string, string, error) {
	// Generate unique key for the file
	fileExt := filepath.Ext(data.FileName)
	uniqueID := uuid.New().String()
	s3Key := fmt.Sprintf("documents/%s%s", uniqueID, fileExt)

	// Upload to S3
	_, err := s.uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(s3Key),
		Body:        bytes.NewReader(data.FileContent),
		ContentType: aws.String(data.ContentType),
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Generate S3 URL
	s3URL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, s3Key)

	return s3Key, s3URL, nil
}

// GetFileURL generates a pre-signed URL for downloading a file
func (s *S3Service) GetFileURL(s3Key string, expiration time.Duration) (string, error) {
	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s3Key),
	})

	url, err := req.Presign(expiration)
	if err != nil {
		return "", fmt.Errorf("failed to generate pre-signed URL: %w", err)
	}

	return url, nil
}

// DeleteFile deletes a file from S3
func (s *S3Service) DeleteFile(s3Key string) error {
	_, err := s.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

// GetFile downloads a file from S3
func (s *S3Service) GetFile(s3Key string) ([]byte, error) {
	buff := &aws.WriteAtBuffer{}
	downloader := s3manager.NewDownloaderWithClient(s.client)

	_, err := downloader.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file from S3: %w", err)
	}

	return buff.Bytes(), nil
}
