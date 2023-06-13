package notifications

import (
	"log"
)

type MockVerificationCodeManager struct {
	AccountCodes map[string]string
	PassCodes    map[string]string
}

func NewMockVerificationCodeManager() *MockVerificationCodeManager {
	return &MockVerificationCodeManager{
		AccountCodes: make(map[string]string),
		PassCodes:    make(map[string]string),
	}
}

func (m *MockVerificationCodeManager) SendEmailToVerifyAccount(code string, email string) error {
	m.AccountCodes[email] = code
	log.Println("Send verification account code:", code, "to:", email)
	return nil
}

func (m *MockVerificationCodeManager) SendEmailToRecoverPassword(code string, email string) error {
	m.PassCodes[email] = code
	log.Println("Send recover password code:", code, "to:", email)
	return nil
}
