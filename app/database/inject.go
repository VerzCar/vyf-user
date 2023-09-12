package database

import (
	"github.com/VerzCar/vyf-lib-logger"
	model2 "github.com/VerzCar/vyf-user/api/model"
	"gorm.io/gorm"
)

func InjectInitData(db *gorm.DB, log logger.Logger) {
	user := &model2.User{
		IdentityID: "2df20ca2-c5b5-4504-ab3b-c800ab73a623",
		FirstName:  "Martin",
		LastName:   "Hammer",
		Gender:     model2.GenderMan,
		Locale:     &model2.Locale{ID: 83},
		Address: &model2.Address{
			Address:    "Badestr. 4",
			City:       "Laufen",
			PostalCode: "83410",
			Country:    &model2.Country{ID: 65},
		},
	}

	if err := db.Create(user).Error; err != nil {
		log.Fatalf("Creation failed: %s", err)
	}

	log.Info("initial data injected.")
}
