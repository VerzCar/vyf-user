package model

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	_ "github.com/go-playground/validator/v10"
	"io"
	"strconv"
	"time"
)

type User struct {
	ID         int64         `json:"id" gorm:"primary_key;index;"`
	IdentityID string        `json:"identityId" gorm:"type:varchar(50);unique;not null"`
	FirstName  string        `json:"firstName" gorm:"type:varchar(50)"`
	LastName   string        `json:"LastName" gorm:"type:varchar(50)"`
	Gender     Gender        `json:"gender" gorm:"type:gender;not null;default:X"`
	LocaleID   sql.NullInt64 `json:"LocaleId"`
	Locale     *Locale       `json:"locale" gorm:"constraint:OnDelete:RESTRICT;"`
	AddressID  sql.NullInt64 `json:"addressId"`
	Address    *Address      `json:"address" gorm:"constraint:OnDelete:CASCADE;"`
	ContactID  sql.NullInt64 `json:"contactId"`
	Contact    *Contact      `json:"contact" gorm:"constraint:OnDelete:CASCADE;"`
	AvatarUrl  string        `json:"avatarUrl"`
	CompanyID  *int64        `json:"companyId"`
	CreatedAt  time.Time     `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt  time.Time     `json:"updatedAt" gorm:"autoUpdateTime;"`
}

type UserCreateInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gt=0"`
}

type UserUpdateInput struct {
	FirstName *string         `json:"firstName" validate:"omitempty,gt=0,lte=50"`
	LastName  *string         `json:"lastName" validate:"omitempty,gt=0,lte=50"`
	Gender    *Gender         `json:"gender"`
	Locale    *string         `json:"locale" validate:"omitempty,bcp47_language_tag"`
	AvatarURL *string         `json:"avatarUrl" validate:"omitempty,url"`
	Address   *AddressInput   `json:"address" validate:"omitempty"`
	Contact   *ContactInputXs `json:"contact" validate:"omitempty"`
}

type UserResetPasswordInput struct {
	ResetPasswordToken string `json:"resetPasswordToken"`
	NewPassword        string `json:"newPassword"`
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

func (e *Gender) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Gender(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Gender", str)
	}
	return nil
}

func (e Gender) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
