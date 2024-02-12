package model

import (
	_ "github.com/go-playground/validator/v10"
	"time"
)

type Address struct {
	CreatedAt  time.Time `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"autoUpdateTime;"`
	Country    *Country  `json:"country" gorm:"constraint:OnDelete:RESTRICT;"`
	Address    string    `json:"address" gorm:"type:varchar(100);not null;check:address != '';"`
	City       string    `json:"city" gorm:"type:varchar(80);not null;check:city != '';"`
	PostalCode string    `json:"postalCode" gorm:"type:varchar(15);not null;check:postal_code != '';"`
	ID         int64     `json:"id" gorm:"primary_key;index;"`
	CountryID  int64     `json:"countryId" gorm:"not null;"`
}

type AddressRequest struct {
	Address          string `json:"address" validate:"required,gt=0,lte=100"`
	City             string `json:"city" validate:"required,gt=0,lte=80"`
	PostalCode       string `json:"postalCode" validate:"required,numeric,gt=0,lte=15"`
	CountryAlphaCode string `json:"countryAlphaCode" validate:"required,iso3166_1_alpha2"`
}
