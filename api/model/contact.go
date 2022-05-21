package model

import (
	"database/sql"
	_ "github.com/go-playground/validator/v10"
	"time"
)

type Contact struct {
	ID                    int64         `json:"id" gorm:"primary_key;index;"`
	Email                 string        `json:"email" gorm:"type:varchar(150);not null;check:email != '';index;"`
	PhoneNumber           string        `json:"phoneNumber" gorm:"type:varchar(20);not null;"`
	PhoneNumberCountryID  int64         `json:"phoneNumberCountryId" gorm:"not null;"`
	PhoneNumberCountry    *Country      `json:"phoneNumberCountry" gorm:"constraint:OnDelete:RESTRICT;"`
	PhoneNumber2          string        `json:"phoneNumber2" gorm:"type:varchar(20);"`
	PhoneNumber2CountryID sql.NullInt64 `json:"phoneNumber2CountryId"`
	PhoneNumber2Country   *Country      `json:"phoneNumber2Country" gorm:"constraint:OnDelete:RESTRICT;"`
	Web                   string        `json:"web" gorm:"type:text;"`
	CreatedAt             time.Time     `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt             time.Time     `json:"updatedAt" gorm:"autoUpdateTime;"`
}

type ContactInput struct {
	Email                        string  `json:"email" validate:"required,email"`
	PhoneNumber                  string  `json:"phoneNumber" validate:"required,numeric,lte=20"`
	PhoneNumberCountryAlphaCode  string  `json:"phoneNumberCountryAlphaCode" validate:"required,iso3166_1_alpha2"`
	PhoneNumber2                 *string `json:"phoneNumber2" validate:"omitempty,numeric,lte=20"`
	PhoneNumber2CountryAlphaCode *string `json:"phoneNumber2CountryAlphaCode" validate:"omitempty,iso3166_1_alpha2"`
	Web                          *string `json:"web" validate:"omitempty,url"`
}

type ContactInputXs struct {
	PhoneNumber                  string  `json:"phoneNumber" validate:"required,numeric,lte=20"`
	PhoneNumberCountryAlphaCode  string  `json:"phoneNumberCountryAlphaCode" validate:"required,iso3166_1_alpha2"`
	PhoneNumber2                 *string `json:"phoneNumber2" validate:"omitempty,numeric,lte=20"`
	PhoneNumber2CountryAlphaCode *string `json:"phoneNumber2CountryAlphaCode" validate:"omitempty,iso3166_1_alpha2"`
	Web                          *string `json:"web" validate:"omitempty,url"`
}
