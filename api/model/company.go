package model

import (
	_ "github.com/go-playground/validator/v10"
	"time"
)

type Company struct {
	ID           int64     `json:"id" gorm:"primary_key;index;"`
	Name         string    `json:"name" gorm:"type:varchar(100);not null;check:name != '';"`
	AddressID    int64     `json:"addressId" gorm:"not null;"`
	Address      *Address  `json:"address" gorm:"constraint:OnDelete:CASCADE;"`
	ContactID    int64     `json:"contactId" gorm:"not null;"`
	Contact      *Contact  `json:"contact" gorm:"constraint:OnDelete:CASCADE;"`
	OwnerID      int64     `json:"ownerId" gorm:"not null;"`
	Owner        *User     `json:"owner" gorm:"foreignKey:OwnerID;constraint:OnDelete:RESTRICT;"`
	Users        []*User   `json:"users" gorm:"constraint:OnDelete:RESTRICT;"`
	BrandLogoUrl string    `json:"brandLogoUrl" gorm:"type:text;"`
	TaxID        string    `json:"taxId" gorm:"type:varchar(20);"`
	IsVerified   bool      `json:"isVerified" gorm:"not null;default:false;"`
	CreatedAt    time.Time `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"autoUpdateTime;"`
}

type CompanyCreateInput struct {
	Name      string        `json:"name" validate:"required,gt=0,lte=100"`
	Address   *AddressInput `json:"address" validate:"required"`
	Contact   *ContactInput `json:"contact" validate:"required"`
	TaxID     string        `json:"taxId" validate:"required,gt=0,lte=20"`
	PaymentId string        `json:"paymentId" validate:"required,uuid4"`
}

type CompanyVerifyInput struct {
	VerificationToken string `json:"verificationToken"`
}
