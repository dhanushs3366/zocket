package services

import (
	"log"

	"github.com/dhanushs3366/zocket/models"
)

func (s *Store) CreateTask(task *models.Task, userID uint) (*models.Task, error) {
	err := s.db.Create(task).Error
	if err != nil {
		log.Println("Couldn't insert new task")
		return nil, err
	}

	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		log.Println("User not found")
		return nil, err
	}

	err = s.db.Model(&user).Association("Tasks").Append(task)
	if err != nil {
		log.Println("Couldn't assign task to user")
		return nil, err
	}

	if err := s.db.Preload("Users").First(&task, task.ID).Error; err != nil {
		log.Println("Couldn't preload users for task")
		return nil, err
	}

	return task, nil
}

func (s *Store) GetTask(ID string) (*models.Task, error) {
	task := new(models.Task)
	if err := s.db.Preload("Users").First(task, "id = ?", ID).Error; err != nil {
		log.Println("Task not found")
		return nil, err
	}
	return task, nil
}

func (s *Store) GetAllTasks(userID string) ([]models.Task, error) {
	var tasks []models.Task
	err := s.db.Preload("Users").
		Joins("JOIN user_tasks ON user_tasks.task_id = tasks.id").
		Where("user_tasks.user_id = ?", userID).
		Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *Store) UpdateTask(task *models.Task) error {
	return s.db.Model(&models.Task{}).Where("id = ?", task.ID).Updates(task).Error
}

func (s *Store) DeleteTask(taskID string) error {
	if err := s.db.Where("id=?", taskID).Delete(&models.Task{}).Error; err != nil {
		log.Printf("couldnt delete task id:%s", taskID)
		return err
	}

	return nil
}
