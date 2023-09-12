package model

import (
	"github.com/VerzCar/vyf-user/utils"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type Profile struct {
	ID                     int64     `json:"id" gorm:"primary_key;index;"`
	Bio                    string    `json:"bio" gorm:"type:varchar(200);not null;"`
	WhyVoteMe              string    `json:"whyVoteMe" gorm:"type:varchar(50);not null;"`
	ImageSrc               string    `json:"imageSrc" gorm:"type:text;not null;"`
	ImagePlaceholderColors string    `json:"imagePlaceholderColors" gorm:"type:varchar(15);not null;"`
	CreatedAt              time.Time `json:"createdAt" gorm:"autoCreateTime;default:current_timestamp;"`
	UpdatedAt              time.Time `json:"updatedAt" gorm:"autoUpdateTime;default:current_timestamp;"`
}

type ProfileInput struct {
	Bio       *string `json:"bio" validate:"omitempty,gt=0,lte=200"`
	WhyVoteMe *string `json:"whyVoteMe" validate:"omitempty,gt=0,lte=50"`
	ImageSrc  *string `json:"imageSrc" validate:"omitempty,url"`
}

func (p *Profile) BeforeCreate(tx *gorm.DB) (err error) {
	gradients := utils.GradientPairList()
	rand.Seed(time.Now().UnixNano())
	randomColor := gradients[rand.Intn(len(gradients))]
	p.ImagePlaceholderColors = randomColor
	return
}
