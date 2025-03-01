package services

import (
	"fmt"
	"log"
	"os"

	"github.com/dhanushs3366/zocket/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func Init() (*Store, error) {
	DB_HOST := os.Getenv("DB_HOST")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")
	DB_PORT := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()

	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	log.Printf("connected to db successfully")

	store := Store{db: db}

	if err := store.migrate(); err != nil {
		log.Printf("couldnt populate tables %s", err.Error())
		return nil, err
	}
	log.Printf("Migration successfull")
	return &store, nil
}

func (s *Store) Close() error {
	sqlDB, err := s.db.DB()

	if err != nil {
		log.Printf("Error accessing sql DB interface %v", err.Error())
		return err
	}

	return sqlDB.Close()
}

func (s *Store) migrate() error {
	return s.db.AutoMigrate(
		&models.User{},
		&models.Task{},
	)
}
