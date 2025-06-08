package handlers

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"guitar_processor/cmd/server/request"
	"guitar_processor/cmd/server/utils"
	"guitar_processor/internal/repository"
)

type AuthHandler struct {
	as   utils.AuthService
	repo *repository.UserRepository
}

type AuthResponse struct {
	Token string `json:"token"`
}

// Handle godoc
// @Summary authorization
// @Description authorization and token creation
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.AuthRequest true "Auth data"
// @Success 200 {object} AuthResponse "token"
// @Failure 400 {object} ErrorResponse "error"
// @Failure 401 {object} ErrorResponse "error"
// @Router /auth [post]
func (ah *AuthHandler) Handle(c *fiber.Ctx) error {
	var req request.AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(ErrorResponse{Error: "invalid request"})
	}

	user, err := ah.repo.GetUserByLogin(req.Login)
	if err != nil {
		return c.Status(400).JSON(ErrorResponse{Error: "invalid request"})
	}

	if ah.checkPasswordHash(req.Password, user.Password) {
		token, _ := ah.as.GenerateToken(user.Login)
		return c.JSON(AuthResponse{Token: token})
	}

	return c.Status(401).JSON(ErrorResponse{Error: "unauthorized"})
}

func (ah *AuthHandler) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func NewAuthHandler(repo *repository.UserRepository, as utils.AuthService) *AuthHandler {
	return &AuthHandler{repo: repo, as: as}
}
