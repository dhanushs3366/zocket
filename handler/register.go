package handler

import (
	"net/http"

	"github.com/dhanushs3366/zocket/models"
	"github.com/dhanushs3366/zocket/services"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) register(c *fiber.Ctx) error {

	// *Requires -> username, password,email.
	user := new(models.User)
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"err":     err.Error(),
		})
	}
	hashedPassword, err := services.HashPassword(user.Password)
	user.Password = hashedPassword

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"err":     err.Error(),
		})
	}

	err = h.store.RegisterUser(user)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"err":     err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"result":  "user registered successfully",
	})
}

func (h *Handler) login(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"err":     err.Error(),
		})
	}

	dbUser, err := h.store.GetUserByUsername(user.Username)

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"err":     err.Error(),
		})
	}

	if dbUser == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"err":     "invalid email or password",
		})
	}

	if !services.ComparePassword(dbUser.Password, user.Password) {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"err":     "invalid email or password",
		})
	}

	// user authorized set cookie
	// generate jwt token
	jwtToken, err := services.GenerateJWTToken(user)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"err":     err.Error(),
		})
	}

	c.Locals("user", dbUser)
	c.Cookie(&fiber.Cookie{
		Name:  "auth_token",
		Value: jwtToken,
	})

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"result":  "user logged in",
	})
}
