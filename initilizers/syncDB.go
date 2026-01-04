package initilizers

import (
	"jwt/models"
)

func SyncDB() {
	DB.AutoMigrate(&models.User{})
}
