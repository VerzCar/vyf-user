package model

import (
	"database/sql"
	"database/sql/driver"
	_ "github.com/go-playground/validator/v10"
	"time"
)

type User struct {
	ID         int64         `json:"id" gorm:"primary_key;"`
	IdentityID string        `json:"identityId" gorm:"type:varchar(50);unique;not null"`
	Username   string        `json:"username" gorm:"type:varchar(40);unique;not null"`
	FirstName  string        `json:"firstName" gorm:"type:varchar(50)"`
	LastName   string        `json:"lastName" gorm:"type:varchar(50)"`
	Gender     Gender        `json:"gender" gorm:"type:gender;not null;default:X"`
	LocaleID   sql.NullInt64 `json:"LocaleId"`
	Locale     *Locale       `json:"locale" gorm:"constraint:OnDelete:RESTRICT;"`
	AddressID  sql.NullInt64 `json:"addressId"`
	Address    *Address      `json:"address" gorm:"constraint:OnDelete:CASCADE;"`
	ContactID  sql.NullInt64 `json:"contactId"`
	Contact    *Contact      `json:"contact" gorm:"constraint:OnDelete:CASCADE;"`
	ProfileID  sql.NullInt64 `json:"profileId"`
	Profile    *Profile      `json:"profile" gorm:"constraint:OnDelete:CASCADE;"`
	CreatedAt  time.Time     `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt  time.Time     `json:"updatedAt" gorm:"autoUpdateTime;"`
}

type UserResponse struct {
	ID         int64     `json:"id"`
	IdentityID string    `json:"identityId"`
	Username   string    `json:"username"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"LastName"`
	Gender     Gender    `json:"gender"`
	Locale     *Locale   `json:"locale,omitempty"`
	Address    *Address  `json:"address,omitempty"`
	Contact    *Contact  `json:"contact,omitempty"`
	Profile    *Profile  `json:"profile,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type UserRequest struct {
	IdentityID *string `json:"identityId,omitempty" validate:"omitempty,lte=50"`
}

type UserUpdateRequest struct {
	FirstName *string         `json:"firstName,omitempty" validate:"omitempty,gt=0,lte=50"`
	LastName  *string         `json:"lastName,omitempty" validate:"omitempty,gt=0,lte=50"`
	Username  *string         `json:"username,omitempty" validate:"omitempty,gt=0,lte=40"`
	Gender    *Gender         `json:"gender,omitempty" validate:"omitempty"`
	Locale    *string         `json:"locale,omitempty" validate:"omitempty,bcp47_language_tag"`
	Address   *AddressRequest `json:"address,omitempty" validate:"omitempty"`
	Contact   *ContactRequest `json:"contact,omitempty" validate:"omitempty"`
	Profile   *ProfileRequest `json:"profile,omitempty" validate:"omitempty"`
}

type Gender string

const (
	GenderX     Gender = "X"
	GenderMan   Gender = "MAN"
	GenderWomen Gender = "WOMEN"
)

var AllGender = []Gender{
	GenderX,
	GenderMan,
	GenderWomen,
}

func (e *Gender) Scan(value interface{}) error {
	*e = Gender(value.(string))
	return nil
}

func (e Gender) Value() (driver.Value, error) {
	return string(e), nil
}

func (e Gender) IsValid() bool {
	switch e {
	case GenderX, GenderMan, GenderWomen:
		return true
	}
	return false
}

func (e Gender) String() string {
	return string(e)
}
