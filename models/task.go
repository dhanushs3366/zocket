package models

import (
	"gorm.io/gorm"
)

type TaskStatus string

const (
	Todo       TaskStatus = "todo"
	InProgress TaskStatus = "pending"
	Completed  TaskStatus = "completed"
)

type Task struct {
	gorm.Model
	Title       string     `gorm:"size:255;not null" json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	Users       []User     `gorm:"many2many:user_tasks;" json:"users"`
}
