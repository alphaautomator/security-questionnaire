package models

import (
	"security-questionnaire/pkg/models"
)

// Result represents a questionnaire result stored in the database
type Result struct {
	models.BaseModel
	QuestionnaireID string `gorm:"column:questionnaire_id;not null;index" json:"questionnaire_id"`
	Data            string `gorm:"column:data;type:jsonb" json:"data"`
	Status          string `gorm:"column:status;not null;default:'pending'" json:"status"`
	Score           *int   `gorm:"column:score" json:"score,omitempty"`
	CompletedAt     *int64 `gorm:"column:completed_at" json:"completed_at,omitempty"`
}

// TableName specifies the table name for the Result model
func (Result) TableName() string {
	return "results"
}
