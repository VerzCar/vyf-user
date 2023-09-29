package model

import (
	"time"
)

type Profile struct {
	ID                     int64     `json:"id" gorm:"primary_key;index;"`
	Bio                    string    `json:"bio" gorm:"type:varchar(1500);not null;"`
	WhyVoteMe              string    `json:"whyVoteMe" gorm:"type:varchar(250);not null;"`
	ImageSrc               string    `json:"imageSrc" gorm:"type:text;not null;"`
	ImagePlaceholderColors string    `json:"imagePlaceholderColors" gorm:"type:varchar(15);not null;"`
	CreatedAt              time.Time `json:"createdAt" gorm:"autoCreateTime;default:current_timestamp;"`
	UpdatedAt              time.Time `json:"updatedAt" gorm:"autoUpdateTime;default:current_timestamp;"`
}

type ProfileRequest struct {
	Bio       *string `json:"bio" validate:"omitempty,gt=0,lte=1500"`
	WhyVoteMe *string `json:"whyVoteMe" validate:"omitempty,gt=0,lte=250"`
	ImageSrc  *string `json:"imageSrc" validate:"omitempty,url"`
}
