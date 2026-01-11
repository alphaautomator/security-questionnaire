package models

import (
	"security-questionnaire/pkg/models"
)

// Document represents a document stored in S3 with metadata in the database
type Document struct {
	models.BaseModel
	FileName    string `gorm:"column:file_name;not null" json:"file_name"`
	FileSize    int64  `gorm:"column:file_size;not null" json:"file_size"`
	ContentType string `gorm:"column:content_type;not null" json:"content_type"`
	S3Bucket    string `gorm:"column:s3_bucket;not null" json:"s3_bucket"`
	S3Key       string `gorm:"column:s3_key;not null;uniqueIndex" json:"s3_key"`
	S3URL       string `gorm:"column:s3_url;not null" json:"s3_url"`
	Description string `gorm:"column:description;type:text" json:"description,omitempty"`
	Tags        string `gorm:"column:tags;type:text" json:"tags,omitempty"`
}

// TableName specifies the table name for the Document model
func (Document) TableName() string {
	return "documents"
}
