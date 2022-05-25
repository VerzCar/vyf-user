package model

import (
	"time"
)

type Bio struct {
	ID          int64     `json:"id" gorm:"primary_key;index;"`
	Description string    `json:"description" gorm:"type:varchar(600);not null;"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime;default:current_timestamp;"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime;default:current_timestamp;"`
}

type BioInput struct {
	Description string `json:"description" validate:"required"`
}
