package model

import (
	"time"
)

type Locale struct {
	CreatedAt    time.Time `json:"createdAt" gorm:"autoCreateTime;default:current_timestamp;"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"autoUpdateTime;default:current_timestamp;"`
	Locale       string    `json:"locale" gorm:"type:varchar(29);not null;"`
	LanguageCode string    `json:"languageCode" gorm:"type:varchar(9);"`
	LcidString   string    `json:"lcidString" gorm:"type:varchar(9);"`
	ID           int64     `json:"id" gorm:"primary_key;index;"`
}
