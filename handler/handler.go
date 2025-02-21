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
	userGroup.Use(services.ValidateJWT)

	// generic reqs
	h.router.Post("/register", h.register)
	h.router.Get("/login", h.login)

	// user reqs
	userGroup.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("value=" + c.Cookies("auth_token"))
	})

	return &h, nil
}

func (h *Handler) Run(port uint) error {
	return h.router.Listen(fmt.Sprintf(":%d", port))
}
