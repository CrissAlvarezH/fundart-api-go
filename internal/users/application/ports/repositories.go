package ports

import (
	"errors"
	users "github.com/CrissAlvarezH/fundart-api/internal/users/domain"
)

var (
	UserDoesNotExists    = errors.New("user does not exists")
	AddressDoesNotExists = errors.New("address does not exists")
	InvalidCredentials   = errors.New("invalid credentials")
)

type UserRepository interface {
	List(filters map[string]string, limit int, offset int) ([]users.User, int)
	GetByID(ID users.UserID) (users.User, bool)
	GetByEmail(email string) (users.User, bool)
	GetPassword(ID users.UserID) (string, bool)
	Add(
		name string, email string, password string, phone string, isActive bool, scopes []users.ScopeName,
	) (users.User, error)
	Update(
		ID users.UserID, name string, email string, phone string, scopes []users.ScopeName,
	) (users.User, error)
	ChangePassword(ID users.UserID, newPassword string) error
	Deactivate(ID users.UserID) error
	Activate(ID users.UserID) bool

	SaveAccountVerificationCode(ID users.UserID, code string) error
	ValidateAccountVerificationCode(ID users.UserID, code string) bool

	SaveRecoveryPasswordCode(ID users.UserID, code string) error
	ValidateRecoveryPasswordCode(ID users.UserID, code string) bool
}

type AddressRepository interface {
	List(ID users.UserID) []users.Address
	Add(
		userID users.UserID, department string, city string, address string,
		receiverPhone string, receiverName string,
	) (users.Address, error)
	Update(
		ID users.AddressID, department string, city string, address string,
		receiverPhone string, receiverName string,
	) (users.Address, error)
	Delete(ID users.AddressID) error
}
