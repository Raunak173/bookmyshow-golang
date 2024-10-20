package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func ConnectToDB() {
	var err error
	dsn := os.Getenv("DB_URL")
	Db, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatal("Error in connection to DB")
	}
}
