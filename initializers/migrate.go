package initializers

import "github.com/raunak173/bms-go/models"

func SyncDB() {
	Db.AutoMigrate(
		&models.Movie{},
		&models.User{},
		&models.Venue{},
		&models.ShowTime{},
		&models.Seat{},
		&models.Order{},
	)
}
