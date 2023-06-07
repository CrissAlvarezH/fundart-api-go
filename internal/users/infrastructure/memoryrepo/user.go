package memoryrepo

import (
	"errors"
	users "github.com/CrissAlvarezH/fundart-api/internal/users/domain"
	"time"
)

type MemoryUser struct {
	ID               users.UserID
	Name             string
	Email            string
	Password         string
	Phone            string
	IsActive         bool
	CreatedAt        time.Time
	Addresses        []users.Address
	Scopes           []users.ScopeName
	VerificationCode string
}

func mapToMemoryUser(user users.User, password string) MemoryUser {
	return MemoryUser{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  password,
		Phone:     user.Phone,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		Addresses: user.Addresses,
		Scopes:    user.Scopes,
	}
}

func mapToUser(memoryUser MemoryUser) users.User {
	return users.User{
		ID:        memoryUser.ID,
		Name:      memoryUser.Name,
		Email:     memoryUser.Email,
		Phone:     memoryUser.Phone,
		IsActive:  memoryUser.IsActive,
		CreatedAt: memoryUser.CreatedAt,
		Addresses: memoryUser.Addresses,
		Scopes:    memoryUser.Scopes,
	}
}

type MemoryUserRepository struct {
	users []MemoryUser
}

func NewMemoryUserRepository(users []MemoryUser) *MemoryUserRepository {
	return &MemoryUserRepository{users: users}
}

func (r *MemoryUserRepository) List(
	filters map[string]string, limit int, offset int,
) ([]users.User, int) {
	filtered := make([]MemoryUser, 0, len(r.users))
	if len(filters) == 0 {
		filtered = r.users
	} else {
		for _, v := range r.users {
			name, ok := filters["name"]
			if ok == true && v.Name == name {
				filtered = append(filtered, v)
			}

			email, ok := filters["email"]
			if ok == true && v.Email == email {
				filtered = append(filtered, v)
			}

			phone, ok := filters["phone"]
			if ok == true && v.Phone == phone {
				filtered = append(filtered, v)
			}
		}
	}

	onlyActives := make([]MemoryUser, 0, len(filtered))
	for _, u := range filtered {
		if u.IsActive == true {
			onlyActives = append(onlyActives, u)
		}
	}
	filtered = onlyActives

	data := make([]MemoryUser, 0, len(filtered))
	if len(filtered) >= limit && len(filtered) >= offset {
		data = filtered[offset:limit]
	} else if len(filtered) >= offset && len(filtered) < limit {
		data = filtered[offset:]
	}

	result := make([]users.User, 0, len(data))
	for _, u := range data {
		result = append(result, mapToUser(u))
	}
	return result, len(filtered)
}

func (r *MemoryUserRepository) GetByID(ID users.UserID) (users.User, bool) {
	for _, u := range r.users {
		if u.ID == ID {
			return mapToUser(u), true
		}
	}
	return users.User{}, false
}

func (r *MemoryUserRepository) Add(
	name string, email string, password string, phone string, isActive bool, scopes []users.ScopeName,
) (users.User, error) {
	lastUser := r.users[len(r.users)-1]
	newUser := users.User{
		ID:       lastUser.ID + 1,
		Name:     name,
		Email:    email,
		Phone:    phone,
		IsActive: isActive,
		Scopes:   scopes,
	}
	r.users = append(r.users, mapToMemoryUser(newUser, password))
	return newUser, nil
}

func (r *MemoryUserRepository) Update(
	ID users.UserID, name string, email string, phone string, isActive bool, scopes []users.ScopeName,
) (users.User, error) {
	return users.User{}, nil
}

func (r *MemoryUserRepository) Deactivate(ID users.UserID) error {
	return nil
}

func (r *MemoryUserRepository) Activate(ID users.UserID) bool {
	for i, u := range r.users {
		if u.ID == ID {
			r.users[i].IsActive = true
			return true
		}
	}
	return false
}

func (r *MemoryUserRepository) ListAddresses(ID users.UserID) []users.Address {
	return make([]users.Address, 0)
}

func (r *MemoryUserRepository) AttachAddress(ID users.UserID, addressID users.AddressID) error {
	return nil
}

func (r *MemoryUserRepository) DetachAddress(ID users.UserID, addressID users.AddressID) error {
	return nil
}

func (r *MemoryUserRepository) SaveVerificationCode(ID users.UserID, code string) error {
	found := false
	for i, u := range r.users {
		if u.ID == ID {
			r.users[i].VerificationCode = code
			found = true
			break
		}
	}

	if found == false {
		return errors.New("user not found")
	}

	return nil
}

func (r *MemoryUserRepository) ValidateVerificationCode(ID users.UserID, code string) bool {
	for _, u := range r.users {
		if u.ID == ID {
			return u.VerificationCode == code
		}
	}
	return false
}
