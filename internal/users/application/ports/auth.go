package ports

import (
	"errors"
	users "github.com/CrissAlvarezH/fundart-api/internal/users/domain"
)

var (
	InvalidToken = errors.New("invalid token")
)

type Token struct {
	AccessToken  string
	RefreshToken string
}

type PasswordManager interface {
	Encrypt(rawPassword string) (string, error)
	Verify(rawPassword string, encryptedPassword string) (bool, error)
}

type JWTManager interface {
	Create(user users.User) (Token, error)
	Verify(accessToken string) (users.User, error)
}
