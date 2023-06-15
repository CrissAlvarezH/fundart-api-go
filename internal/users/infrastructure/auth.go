package infrastructure

import (
	"github.com/CrissAlvarezH/fundart-api/internal/users/application/ports"
	users "github.com/CrissAlvarezH/fundart-api/internal/users/domain"
	"strings"
)

type MockPasswordManager struct {
}

func NewMockPasswordManager() *MockPasswordManager {
	return &MockPasswordManager{}
}

func (p *MockPasswordManager) Encrypt(rawPassword string) (string, error) {
	return rawPassword + "_encrypt", nil
}

func (p *MockPasswordManager) Verify(rawPassword string, encryptedPassword string) (bool, error) {
	return encryptedPassword == rawPassword+"_encrypt", nil
}

type MockJWTManager struct {
	userRepo ports.UserRepository
}

func NewMockJWTManager(userRepo ports.UserRepository) *MockJWTManager {
	return &MockJWTManager{userRepo: userRepo}
}

func (m *MockJWTManager) Create(user users.User) (ports.Token, error) {
	return ports.Token{
		AccessToken:  user.Email + "___jwt",
		RefreshToken: user.Email + "__refresh",
	}, nil
}

func (m *MockJWTManager) Verify(accessToken string) (users.User, error) {
	email := strings.Split(accessToken, "___")[0]
	user, ok := m.userRepo.GetByEmail(email)
	if !ok {
		return users.User{}, ports.InvalidToken
	}

	return user, nil
}
