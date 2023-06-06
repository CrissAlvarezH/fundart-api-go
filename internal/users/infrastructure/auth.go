package infrastructure

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
