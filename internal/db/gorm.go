// internal/db/gorm.go

package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"hyberica.io/go/go-api/internal/domain/models"
)

// NewDatabaseConnection cria e retorna uma nova conexão do banco de dados usando GORM.
// Esta função pode ser usada pelo Wire para injeção de dependência.
func NewDatabaseConnection() (*gorm.DB, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	db.AutoMigrate(&models.Tenant{})

	log.Printf("INFO: DB (Gorm) inicializado com sucesso!")
	return db, nil
}
