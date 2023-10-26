package repository

import (
	"database/sql"
	"fmt"
	"github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-user/api/model"
	"github.com/VerzCar/vyf-user/app/config"
	"github.com/VerzCar/vyf-user/app/database"
	"github.com/VerzCar/vyf-user/utils"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"path/filepath"
)

type Storage interface {
	RunMigrationsUp(db *sql.DB) error
	RunMigrationsDown(db *sql.DB) error

	CountryById(id int64) (*model.Country, error)
	CountryByAlpha2(alpha2 string) (*model.Country, error)

	LocaleByLcidString(lcid string) (*model.Locale, error)

	TransformAddressRequest(src *model.AddressRequest, dest *model.Address) error
	TransformContactRequest(src *model.ContactRequest, dest *model.Contact) error

	CreateNewUser(user *model.User) (*model.User, error)
	UpdateUser(user *model.User) (*model.User, error)
	UserById(id int64) (*model.User, error)
	UserByIdentityId(id string) (*model.User, error)
	Users(identityID string) ([]*model.UserPaginated, error)
	UsersFiltered(
		callerIdentityID string,
		username string,
	) ([]*model.UserPaginated, error)
}

type storage struct {
	db     database.Client
	config *config.Config
	log    logger.Logger
}

func NewStorage(
	db database.Client,
	config *config.Config,
	log logger.Logger,
) Storage {
	return &storage{
		db:     db,
		config: config,
		log:    log,
	}
}

func (s *storage) RunMigrationsUp(db *sql.DB) error {
	m, err := createMigrateInstance(db)

	if err != nil {
		return err
	}

	err = m.Up()

	switch err {
	case migrate.ErrNoChange:
		return nil
	}

	return err
}

func (s *storage) RunMigrationsDown(db *sql.DB) error {
	m, err := createMigrateInstance(db)

	if err != nil {
		return err
	}

	err = m.Down()

	switch err {
	case migrate.ErrNoChange:
		return nil
	}

	return err
}

func createMigrateInstance(db *sql.DB) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		return nil, fmt.Errorf("error creating migrations with db instance: %s", err)
	}

	repoMigrationPath := utils.FromBase("repository/migrations/")
	migrationsPath := filepath.Join("file://", repoMigrationPath)

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres", driver,
	)

	if err != nil {
		return nil, err
	}

	return m, nil
}
