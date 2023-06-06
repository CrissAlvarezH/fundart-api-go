package ports

type PasswordManager interface {
	Encrypt(rawPassword string) (string, error)
	Verify(rawPassword string, encryptedPassword string) (bool, error)
}
