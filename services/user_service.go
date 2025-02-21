package services

import (
	"log"

	"github.com/dhanushs3366/zocket/models"
)

func (s *Store) RegisterUser(user *models.User) error {
	// requires hashed password

	result := s.db.Create(user)

	if result.Error != nil {
		log.Printf("err creating user %s", result.Error.Error())
		return result.Error
	}
	return nil
}

func (s *Store) GetUserByUsername(name string) (*models.User, error) {
	var user models.User
	result := s.db.First(&user, "name = ?", name)
	if result.Error != nil {
		log.Println("user not found")
		return nil, result.Error
	}
	return &user, nil
}
