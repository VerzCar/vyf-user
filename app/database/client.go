package database

import (
	"database/sql"
	"gorm.io/gorm"
)

type Client interface {
	DB() (*sql.DB, error)
	Model(value interface{}) (tx *gorm.DB)
	Session(config *gorm.Session) *gorm.DB
	Create(value interface{}) (tx *gorm.DB)
	Preload(query string, args ...interface{}) (tx *gorm.DB)
	Where(query interface{}, args ...interface{}) (tx *gorm.DB)
}
