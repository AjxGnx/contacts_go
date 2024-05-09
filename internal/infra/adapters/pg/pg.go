package pg

import (
	"fmt"
	"sync"

	"github.com/AjxGnx/contacts-go/config"
	"github.com/AjxGnx/contacts-go/internal/domain/models"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	instance *gorm.DB
	once     sync.Once
)

func ConnInstance() *gorm.DB {
	once.Do(func() {
		instance = getConnection()
	})

	return instance
}

func getConnection() *gorm.DB {
	connString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=disable",
		config.Environments().DBHost,
		config.Environments().DBUser,
		config.Environments().DBPass,
		config.Environments().DBName,
		config.Environments().DBPort)

	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err = db.AutoMigrate(models.Contact{}); err != nil {
		log.Fatal(err)
	}

	return db
}
