package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/valentergs/bibli-users/api/models"
)

var users = []models.User{
	models.User{
		FirstName: "Rodrigo",
		LastName:  "Valente",
		Email:     "valentergs@gmail.com",
		Password:  "password",
		Role:      "admin",
		Active:    true,
		Photo:     "https://robohash.org/rodrigo+valente?set=set2",
	},
	models.User{
		FirstName: "Bruce",
		LastName:  "Wayne",
		Email:     "batman@gmail.com",
		Password:  "password",
		Role:      "staff",
		Active:    true,
		Photo:     "https://robohash.org/bruce+wayne?set=set2",
	},
	models.User{
		FirstName: "Peter",
		LastName:  "Parker",
		Email:     "spiderman@gmail.com",
		Password:  "password",
		Role:      "customer",
		Active:    true,
		Photo:     "https://robohash.org/peter+parker?set=set2",
	},
	models.User{
		FirstName: "Edson",
		LastName:  "Arantes do Nascimento",
		Email:     "pele@gmail.com",
		Password:  "password",
		Role:      "user",
		Active:    true,
		Photo:     "https://robohash.org/edson+arantes?set=set2",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}

	err = db.Debug().AutoMigrate(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}

}
