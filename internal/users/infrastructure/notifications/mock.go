package notifications

import (
	"log"
)

type MockVerificationCodeManager struct {
}

func NewMockVerificationCodeManager() *MockVerificationCodeManager {
	return &MockVerificationCodeManager{}
}

func (m *MockVerificationCodeManager) SendEmailToVerifyAccount(code string, email string) error {
	log.Println("Send verification account code:", code, "to:", email)
	return nil
}

func (m *MockVerificationCodeManager) SendEmailToRecoverPassword(code string, email string) error {
	log.Println("Send recover password code:", code, "to:", email)
	return nil
}
