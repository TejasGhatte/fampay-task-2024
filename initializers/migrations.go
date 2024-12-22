package initializers

import (
	"log"
	"fmt"

	"github.com/TejasGhatte/fampay-task-2024/models"

)

func DBMigrate() {
	if DB == nil {
		fmt.Println("Database not connected")
	}

    err := DB.AutoMigrate(
		&models.Video{},
	)

	if err != nil {
		log.Fatal("Failed to run migrations")
	}
}