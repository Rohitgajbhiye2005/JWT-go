package initilizers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {

	dsn := os.Getenv("DB_URL")

	if dsn == "" {
		log.Println("Database url is Empty")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: false,
	})

	if err != nil {
		log.Fatal("DB error ", err)
	}

	DB = db

	log.Println("Database is Connected")
}
