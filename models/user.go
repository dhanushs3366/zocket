package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"size:100; not null" json:"name"`
	Username string `gorm:"unique;size:100;not null" json:"username"`
	Email    string `gorm:"size:255;unique;not null" json:"email"`
	Password string `gorm:"size:255;not null" json:"password"`  // hashed password
	Tasks    []Task `gorm:"many2many:user_tasks;" json:"tasks"` // many2many because a group of user can have same task for like collabaration
}
