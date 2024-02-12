package model

import (
	"database/sql"
	_ "github.com/go-playground/validator/v10"
	"time"
)

type Contact struct {
	CreatedAt             time.Time     `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt             time.Time     `json:"updatedAt" gorm:"autoUpdateTime;"`
	PhoneNumberCountry    *Country      `json:"phoneNumberCountry" gorm:"constraint:OnDelete:RESTRICT;"`
	PhoneNumber2Country   *Country      `json:"phoneNumber2Country" gorm:"constraint:OnDelete:RESTRICT;"`
	Email                 string        `json:"email" gorm:"type:varchar(150);not null;check:email != '';index;"`
	PhoneNumber           string        `json:"phoneNumber" gorm:"type:varchar(20);not null;"`
	PhoneNumber2          string        `json:"phoneNumber2" gorm:"type:varchar(20);"`
	Web                   string        `json:"web" gorm:"type:text;"`
	PhoneNumber2CountryID sql.NullInt64 `json:"phoneNumber2CountryId"`
	ID                    int64         `json:"id" gorm:"primary_key;index;"`
	PhoneNumberCountryID  int64         `json:"phoneNumberCountryId" gorm:"not null;"`
}

type ContactRequest struct {
	PhoneNumber2                 *string `json:"phoneNumber2" validate:"omitempty,numeric,lte=20"`
	PhoneNumber2CountryAlphaCode *string `json:"phoneNumber2CountryAlphaCode" validate:"omitempty,iso3166_1_alpha2"`
	Web                          *string `json:"web" validate:"omitempty,url"`
	Email                        string  `json:"email" validate:"required,email"`
	PhoneNumber                  string  `json:"phoneNumber" validate:"required,numeric,lte=20"`
	PhoneNumberCountryAlphaCode  string  `json:"phoneNumberCountryAlphaCode" validate:"required,iso3166_1_alpha2"`
}
