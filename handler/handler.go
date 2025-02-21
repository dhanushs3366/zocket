package handler

import (
	"fmt"
	"os"

	"github.com/dhanushs3366/zocket/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
)

type Handler struct {
	router *fiber.App
	store  *services.Store
}

func Init() (*Handler, error) {
	// get db initiated etc
	store, err := services.Init()

	if err != nil {
		return nil, err
	}

	h := Handler{
		router: fiber.New(),
		store:  store,
	}

	// Groups
	userGroup := h.router.Group("users")

	// middlewares
	h.router.Use(encryptcookie.New(encryptcookie.Config{
		Key: os.Getenv("COOKIE_SECRET"),
	}))
	userGroup.Use(services.ValidateJWT(h.store))

	// generic reqs
	h.router.Post("/register", h.register)
	h.router.Get("/login", h.login)

	// user reqs
	userGroup.Post("/tasks", h.createTask)
	userGroup.Get("/tasks/:id", h.GetTask)
	userGroup.Get("/tasks", h.GetAllTasks)
	userGroup.Patch("/tasks", h.UpdateTask)
	userGroup.Delete("/tasks", h.DeleteTask)

	// @TODO: add a feature for group collaboration where a user can add other users to their tasks

	return &h, nil
}

func (h *Handler) Run(port uint) error {
	return h.router.Listen(fmt.Sprintf(":%d", port))
}
