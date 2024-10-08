package database

import (
	"fmt"
	"github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-user/app/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
)

// Connect to the database with configured dsn.
// If successful, the gorm.DB connection will be returned, otherwise
// an error is written and os.exit will be executed.
func Connect(log logger.Logger, conf *config.Config) Client {
	log.Infof("Connect to database...")

	db, err := gorm.Open(
		postgres.Open(dsn(conf)),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		},
	)

	if err != nil {
		log.Fatalf("Connect to database DSN: %s failed: %s", dsn(conf), err)
	}

	log.Infof(
		"Connection to database successfully established via: %s",
		strings.ReplaceAll(dsn(conf), conf.Db.Password, ""),
	)

	return db
}

// dsn construct the database dsn from the configuration.
func dsn(conf *config.Config) string {
	sslMode := "disable"
	if conf.Environment == config.EnvironmentProd {
		sslMode = "require"
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		conf.Db.Host,
		conf.Db.Port,
		conf.Db.User,
		conf.Db.Name,
		conf.Db.Password,
		sslMode,
	)
}
