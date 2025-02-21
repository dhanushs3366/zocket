package handler

import (
	"net/http"
	"strconv"

	"github.com/dhanushs3366/zocket/models"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) createTask(c *fiber.Ctx) error {
	task := new(models.Task)
	user := c.Locals("user").(*models.User)

	if err := c.BodyParser(task); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"err":     err.Error(),
		})
	}
	createdTask, err := h.store.CreateTask(task, user.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"err":     err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"result":  createdTask,
	})
}

func (h *Handler) GetTask(c *fiber.Ctx) error {
	taskID := c.Params("id", "")
	if taskID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"err":     "task_id is missing",
		})
	}

	task, err := h.store.GetTask(taskID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"err":     err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"result":  task,
	})
}

func (h *Handler) GetAllTasks(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	tasks, err := h.store.GetAllTasks(strconv.Itoa(int(user.ID)))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"err":     err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"tasks":   tasks,
	})
}

func (h *Handler) UpdateTask(c *fiber.Ctx) error {

	// remove this if you receive proper task struct from FE
	taskIDStr := c.Query("id")
	taskID32, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"err":     "invalid task ID",
		})
	}

	taskID := uint(taskID32)

	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"err":     err.Error(),
		})
	}
	task.ID = taskID
	if err := h.store.UpdateTask(&task); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"err":     err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Task updated successfully",
	})
}

func (h *Handler) DeleteTask(c *fiber.Ctx) error {
	taskID := c.Query("id")
	if taskID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"err":     "task_id is missing",
		})
	}

	if err := h.store.DeleteTask(taskID); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"err":     err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Task deleted successfully",
	})
}
