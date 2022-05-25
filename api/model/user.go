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
	Username   string        `json:"username" gorm:"type:varchar(40);unique;not null"`
	FirstName  string        `json:"firstName" gorm:"type:varchar(50)"`
	LastName   string        `json:"LastName" gorm:"type:varchar(50)"`
	Gender     Gender        `json:"gender" gorm:"type:gender;not null;default:X"`
	LocaleID   sql.NullInt64 `json:"LocaleId"`
	Locale     *Locale       `json:"locale" gorm:"constraint:OnDelete:RESTRICT;"`
	AddressID  sql.NullInt64 `json:"addressId"`
	Address    *Address      `json:"address" gorm:"constraint:OnDelete:CASCADE;"`
	ContactID  sql.NullInt64 `json:"contactId"`
	Contact    *Contact      `json:"contact" gorm:"constraint:OnDelete:CASCADE;"`
	BioID      sql.NullInt64 `json:"bioId"`
	Bio        *Bio          `json:"bio" gorm:"constraint:OnDelete:CASCADE;"`
	CreatedAt  time.Time     `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt  time.Time     `json:"updatedAt" gorm:"autoUpdateTime;"`
}

type UserUpdateInput struct {
	FirstName *string       `json:"firstName" validate:"omitempty,gt=0,lte=50"`
	LastName  *string       `json:"lastName" validate:"omitempty,gt=0,lte=50"`
	Username  *string       `json:"username" validate:"omitempty,gt=0,lte=40"`
	Gender    *Gender       `json:"gender"`
	Locale    *string       `json:"locale" validate:"omitempty,bcp47_language_tag"`
	Address   *AddressInput `json:"address" validate:"omitempty"`
	Contact   *ContactInput `json:"contact" validate:"omitempty"`
	Bio       *BioInput     `json:"bio" validate:"omitempty"`
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
