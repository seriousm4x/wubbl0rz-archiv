package database

import (
	"fmt"
	"os"
	"strconv"

	"github.com/AgileProggers/archiv-backend-go/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	// Connect DB
	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("host=db user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Europe/Berlin", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&models.Vod{}, &models.Game{}, &models.Creator{}, &models.Clip{}, &models.Emote{}, &models.Settings{})
	if err != nil {
		panic(fmt.Sprint("Unable to auto migrate database:", err))
	}
	DB = db
}

func Close() error {
	db, err := DB.DB()
	if err != nil {
		return fmt.Errorf("gorm.DB get database: %v", err)
	}
	return db.Close()
}
