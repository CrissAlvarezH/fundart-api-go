package notifications

import (
	"github.com/CrissAlvarezH/fundart-api/internal/users/application/ports"
	"log"
)

type MockVerificationCodeManager struct {
}

func NewMockVerificationCodeManager() *MockVerificationCodeManager {
	return &MockVerificationCodeManager{}
}

func (m *MockVerificationCodeManager) Send(code string, provider ports.MessageProvider) error {
	log.Println("Send code:", code, "with provider:", provider)
	return nil
}
