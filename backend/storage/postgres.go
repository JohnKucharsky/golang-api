package storage

import (
	"database/sql"
	"fmt"
	"github.com/JohnKucharsky/golang-api/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func connect() (*gorm.DB, error) {
	errDotEnv := godotenv.Load()
	envDbAddress := os.Getenv("DB_ADDRESS")
	bin, errReadFile := os.ReadFile("/run/secrets/db-password")
	if errReadFile != nil {
		db, _ := sql.Open(
			"postgres",
			envDbAddress,
		)
		return gorm.Open(
			postgres.New(
				postgres.Config{
					Conn: db,
				},
			), &gorm.Config{},
		)
	}

	if errDotEnv != nil {
		db, _ := sql.Open(
			"postgres",
			fmt.Sprintf(
				"postgres://postgres:%s@db:5432/data?sslmode=disable",
				string(bin),
			),
		)

		return gorm.Open(
			postgres.New(
				postgres.Config{
					Conn: db,
				},
			), &gorm.Config{},
		)
	}
	return nil, nil
}

func NewConnection() (*gorm.DB, error) {
	db, err := connect()

	if err != nil {
		return db, err
	}

	return db, nil
}

func MigrateBooks(db *gorm.DB) error {
	err := db.AutoMigrate(&models.Book{})
	return err
}
