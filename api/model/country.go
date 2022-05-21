package model

import (
	"time"
)

type Country struct {
	ID            int64     `json:"id" gorm:"primary_key;index;"`
	Name          string    `json:"name" gorm:"type:varchar(255);not null;"`
	Alpha2        string    `json:"alpha2" gorm:"type:char(2);not null;"`
	Alpha3        string    `json:"alpha3" gorm:"type:char(3);not null;"`
	ContinentCode string    `json:"continentCode" gorm:"type:char(2);not null;"`
	Number        string    `json:"number" gorm:"type:char(3);not null;"`
	FullName      string    `json:"fullName" gorm:"type:varchar(255);not null;"`
	CreatedAt     time.Time `json:"createdAt" gorm:"autoCreateTime;default:current_timestamp;"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"autoUpdateTime;default:current_timestamp;"`
}
