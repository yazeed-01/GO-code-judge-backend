package initializers

import (
	"CfBE/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {

	dsn := os.Getenv("DATABASE_URL")
	var err error

	// Connect to the database using GORM
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v\n", err)
	}

	// Auto-migrate (build) models
	err = DB.AutoMigrate(
		&models.User{},
		&models.Contest{},
		&models.Problem{},
		&models.Language{},
		&models.Submission{},
		&models.Result{},
	)
	if err != nil {
		log.Fatalf("Unable to migrate the database schema: %v\n", err)
	}

	log.Println("Connected to the database and migrated schema successfully")
}
