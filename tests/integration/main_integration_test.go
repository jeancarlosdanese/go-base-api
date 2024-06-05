// tests/integration/main_integration_test.go

package integration_test

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	var err error

	// Configura a conexão com o banco de dados de teste
	dsn, err := getDSN()
	if err != nil {
		log.Fatalf("Failed GetDSN: %v", err)
	}

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	migrationsDir := "file://" + getMigrationPath() // Isso configura o caminho antes de usar.

	// Reverte as migrações
	if err := rollbackMigrations(migrationsDir, dsn); err != nil {
		log.Printf("Failed to rollback migrations: %v", err)
	}

	// Configura e executa as migrações
	if err := applyMigrations("file://"+getMigrationPath(), dsn); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Executa os testes
	code := m.Run()

	// Reverte as migrações
	if err := rollbackMigrations(migrationsDir, dsn); err != nil {
		log.Printf("Failed to rollback migrations: %v", err)
	}

	// Cleanup, se necessário
	os.Exit(code)
}

func applyMigrations(migrationsDir, dsn string) error {
	m, err := migrate.New(
		migrationsDir,
		dsn,
	)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func rollbackMigrations(migrationsDir, dsn string) error {
	m, err := migrate.New(migrationsDir, dsn)
	if err != nil {
		return err
	}
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func getMigrationPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get current directory")
	}
	return filepath.Join(dir, "../../migrations")
}

func getDSN() (string, error) {
	envFile := ".env.test"
	if err := godotenv.Load(envFile); err != nil {
		return "", fmt.Errorf("warning: %s file not found: %w", envFile, err)
	}

	user := url.QueryEscape(os.Getenv("DB_USER"))
	password := url.QueryEscape(os.Getenv("DB_PASSWORD"))
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if user == "" || password == "" || host == "" || port == "" || dbName == "" {
		return "", fmt.Errorf("database environment settings are not fully configured")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&TimeZone=America/Sao_Paulo",
		user, password, host, port, dbName)
	return dsn, nil
}
