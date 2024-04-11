package model

import (
	"database/sql"
	"database/sql/driver"
	_ "github.com/go-playground/validator/v10"
	"time"
)

type User struct {
	UpdatedAt  time.Time     `json:"updatedAt" gorm:"autoUpdateTime;"`
	CreatedAt  time.Time     `json:"createdAt" gorm:"autoCreateTime;"`
	Contact    *Contact      `json:"contact" gorm:"constraint:OnDelete:CASCADE;"`
	Locale     *Locale       `json:"locale" gorm:"constraint:OnDelete:RESTRICT;"`
	Address    *Address      `json:"address" gorm:"constraint:OnDelete:CASCADE;"`
	Profile    *Profile      `json:"profile" gorm:"constraint:OnDelete:CASCADE;"`
	IdentityID string        `json:"identityId" gorm:"type:varchar(50);unique;not null"`
	Username   string        `json:"username" gorm:"type:varchar(40);unique;not null"`
	FirstName  string        `json:"firstName" gorm:"type:varchar(50)"`
	LastName   string        `json:"lastName" gorm:"type:varchar(50)"`
	Gender     Gender        `json:"gender" gorm:"type:gender;not null;default:X"`
	AddressID  sql.NullInt64 `json:"addressId"`
	ProfileID  sql.NullInt64 `json:"profileId"`
	ContactID  sql.NullInt64 `json:"contactId"`
	LocaleID   sql.NullInt64 `json:"LocaleId"`
	ID         int64         `json:"id" gorm:"primary_key;"`
}

type UserResponse struct {
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Locale     *Locale   `json:"locale,omitempty"`
	Address    *Address  `json:"address,omitempty"`
	Contact    *Contact  `json:"contact,omitempty"`
	Profile    *Profile  `json:"profile,omitempty"`
	IdentityID string    `json:"identityId"`
	Username   string    `json:"username"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	Gender     Gender    `json:"gender"`
	ID         int64     `json:"id"`
}

type UserPaginated struct {
	IdentityID      string `json:"identityId"`
	Username        string `json:"username"`
	ProfileImageSrc string `json:"profileImageSrc"`
}

type UserPaginatedResponse struct {
	Profile    *ProfilePaginatedResponse `json:"profile"`
	IdentityID string                    `json:"identityId"`
	Username   string                    `json:"username"`
}

type UserXUriRequest struct {
	IdentityID string `uri:"identityId" validate:"lte=50"`
}

type UserByUriRequest struct {
	Username string `uri:"username" validate:"lte=50"`
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
