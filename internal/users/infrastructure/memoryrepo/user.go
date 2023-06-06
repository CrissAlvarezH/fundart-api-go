package memoryrepo

import (
	users "github.com/CrissAlvarezH/fundart-api/internal/users/domain"
)

type MemoryUserRepository struct {
	users []users.User
}

func NewMemoryUserRepository(users []users.User) *MemoryUserRepository {
	return &MemoryUserRepository{users: users}
}

func (r *MemoryUserRepository) List(
	filters map[string]string, limit int, offset int,
) ([]users.User, int, error) {
	filtered := make([]users.User, 0, len(r.users))
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
	data := make([]users.User, 0, len(filtered))
	if len(filtered) >= limit && len(filtered) >= offset {
		data = filtered[offset:limit]
	} else if len(filtered) >= offset && len(filtered) < limit {
		data = filtered[offset:]
	}
	return data, len(filtered), nil
}

func (r *MemoryUserRepository) GetByID(ID users.UserID) (users.User, bool) {
	for _, u := range r.users {
		if u.ID == ID {
			return u, true
		}
	}
	return users.User{}, false
}

func (r *MemoryUserRepository) Add(
	name string, email string, password string, phone string, isActive bool, scopes []users.ScopeName,
) (users.User, error) {
	return users.User{}, nil
}

func (r *MemoryUserRepository) Update(
	ID users.UserID, name string, email string, phone string, isActive bool, scopes []users.ScopeName,
) (users.User, error) {
	return users.User{}, nil
}

func (r *MemoryUserRepository) Deactivate(ID users.UserID) error {
	return nil
}

func (r *MemoryUserRepository) ListAddress(ID users.UserID) (users.Address, error) {
	return users.Address{}, nil
}

func (r *MemoryUserRepository) AttachAddress(ID users.UserID, addressID users.AddressID) error {
	return nil
}

func (r *MemoryUserRepository) DetachAddress(ID users.UserID, addressID users.AddressID) error {
	return nil
}

func (r *MemoryUserRepository) SaveVerificationCode(ID users.UserID, code string) error {
	return nil
}

func (r *MemoryUserRepository) ValidateVerificationCode(ID users.UserID, code string) (bool, error) {
	return false, nil
}
