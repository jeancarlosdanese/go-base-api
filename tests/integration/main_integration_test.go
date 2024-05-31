package integration_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	var err error

	// Configura a conexão com o banco de dados de teste
	dsn := "postgres://hyberica:hyberica@localhost:5432/go_ead_api_test?sslmode=disable"
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
	if err := runMigrations("file://"+getMigrationPath(), dsn); err != nil {
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

func runMigrations(migrationsDir, dsn string) error {
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
