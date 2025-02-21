package services

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dhanushs3366/zocket/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	HASHING_ROUNDS = 14
	EXPIRY_TIME    = 24 //in hrs
)

type UserClaims struct {
	ID       uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), HASHING_ROUNDS)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func ComparePassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func GenerateJWTToken(user *models.User) (string, error) {
	JWT_SECRET := os.Getenv("JWT_SECRET")
	expirationTime := time.Now().Add(EXPIRY_TIME * time.Hour)

	claims := &UserClaims{
		ID:       user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(JWT_SECRET))

	if err != nil {
		log.Println("couldnt generate jwt token")
		return "", err
	}

	return tokenStr, nil
}

// a middlerware to validate the user
func ValidateJWT(c *fiber.Ctx) error {
	JWT_SECRET := []byte(os.Getenv("JWT_SECRET"))
	cookie := c.Cookies("auth_token")

	if cookie == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "missing auth token"})
	}

	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(cookie, claims, func(t *jwt.Token) (interface{}, error) {
		return JWT_SECRET, nil
	})

	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			status = http.StatusUnauthorized
		}
		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}

	if !token.Valid {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "user unauthorized"})
	}

	return c.Next()
}

func GetUserIDFromToken(c *fiber.Ctx) (uint, error) {
	cookie := c.Cookies("auth_token")
	JWT_SECRET := os.Getenv("JWT_SECRET")

	tokenStr := cookie

	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims.ID, nil
	}

	return 0, errors.New("invalid token")
}
