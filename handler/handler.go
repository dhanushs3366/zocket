package handler

import (
	"fmt"

	"github.com/dhanushs3366/zocket/services"
	"github.com/gofiber/fiber/v2"
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

	// add all methods for the handler

	return &h, nil
}

func (h *Handler) Run(port uint) error {
	return h.router.Listen(fmt.Sprintf(":%d", port))
}
