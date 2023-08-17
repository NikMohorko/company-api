package database

import (
	"company_api/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseInstance struct {
	Db *gorm.DB
}

var DB DatabaseInstance

func Connect() {

	dsn := fmt.Sprintf(
		"host=db user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Europe/London",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Info),
		TranslateError: true, // Get errors from Postgres
	})

	if err != nil {
		log.Fatal("Database connection failed.\n", err)
		os.Exit(2)
	}

	log.Println("Connected to database.")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running database migration.")
	if err = db.AutoMigrate(&models.Company{}, &models.User{}); err != nil {
		log.Println(err.Error())
	}

	DB = DatabaseInstance{
		Db: db,
	}

}
