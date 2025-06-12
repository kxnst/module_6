package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"guitar_processor/internal/entity"
	"guitar_processor/internal/repository"
	"time"
)

var jwtSecret = []byte("super-secret-key")

type Claims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}

type AuthService interface {
	ValidateTokenAndGetUser(tokenStr string) (*entity.User, error)
	GenerateToken(userID string) (string, error)
}
type JWTAuthService struct {
	ur *repository.UserRepository
}

func (a *JWTAuthService) ValidateTokenAndGetUser(tokenStr string) (*entity.User, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	user, err := a.ur.GetUserByLogin(claims.Login)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (a *JWTAuthService) GenerateToken(userID string) (string, error) {
	claims := Claims{
		Login: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func NewAuthService(userRepo *repository.UserRepository) AuthService {
	return &JWTAuthService{ur: userRepo}
}
