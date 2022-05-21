package database

import (
	"gitlab.vecomentman.com/libs/logger"
	model2 "gitlab.vecomentman.com/service/user/api/model"
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

	users := []*model2.User{user}

	company := &model2.Company{
		Name: "Marcos Pizza",
		Address: &model2.Address{
			Address:    "Brienner Str. 49",
			City:       "München",
			PostalCode: "80333",
			Country:    &model2.Country{ID: 65},
		},
		Contact: &model2.Contact{
			Email:               "dev.marcos-pizza@vecomentman.de",
			PhoneNumber:         "08938889733",
			PhoneNumberCountry:  &model2.Country{ID: 65},
			PhoneNumber2:        "08938889755",
			PhoneNumber2Country: &model2.Country{ID: 65},
			Web:                 "https://marcospizza.de",
		},
		TaxID: "DE123456789",
		Owner: user,
		Users: users,
	}

	if err := db.Create(company).Error; err != nil {
		log.Fatalf("Creation failed: %s", err)
	}

	log.Info("initial data injected.")
}
