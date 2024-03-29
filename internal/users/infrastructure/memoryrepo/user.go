package memoryrepo

import (
	"github.com/CrissAlvarezH/fundart-api/internal/users/application/ports"
	users "github.com/CrissAlvarezH/fundart-api/internal/users/domain"
	"time"
)

type MemoryUser struct {
	ID                   users.UserID
	Name                 string
	Email                string
	Password             string
	Phone                string
	IsActive             bool
	CreatedAt            time.Time
	Addresses            []users.Address
	Scopes               []users.ScopeName
	VerificationCode     string
	RecoveryPasswordCode string
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
	Users []MemoryUser
}

func NewMemoryUserRepository(users []MemoryUser) *MemoryUserRepository {
	return &MemoryUserRepository{Users: users}
}

func (r *MemoryUserRepository) List(
	filters map[string]string, limit int, offset int,
) ([]users.User, int) {
	filtered := make([]MemoryUser, 0, len(r.Users))
	if len(filters) == 0 {
		filtered = r.Users
	} else {
		for _, v := range r.Users {
			name, ok := filters["name"]
			if ok && v.Name == name {
				filtered = append(filtered, v)
			}

			email, ok := filters["email"]
			if ok && v.Email == email {
				filtered = append(filtered, v)
			}

			phone, ok := filters["phone"]
			if ok && v.Phone == phone {
				filtered = append(filtered, v)
			}
		}
	}

	onlyActives := make([]MemoryUser, 0, len(filtered))
	for _, u := range filtered {
		if u.IsActive {
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
	for _, u := range r.Users {
		if u.ID == ID {
			if u.IsActive {
				return mapToUser(u), true
			}
			break
		}
	}
	return users.User{}, false
}

func (r *MemoryUserRepository) GetByEmail(email string) (users.User, bool) {
	for _, u := range r.Users {
		if u.Email == email {
			if u.IsActive {
				return mapToUser(u), true
			}
			break
		}
	}
	return users.User{}, false
}

func (r *MemoryUserRepository) GetPassword(ID users.UserID) (string, bool) {
	for _, u := range r.Users {
		if u.ID == ID {
			if u.IsActive {
				return u.Password, true
			}
			break
		}
	}
	return "", false
}

func (r *MemoryUserRepository) Add(
	name string, email string, password string, phone string, isActive bool, scopes []users.ScopeName,
) (users.User, error) {
	lastUserID := users.UserID(0)
	if len(r.Users) > 0 {
		lastUserID = r.Users[len(r.Users)-1].ID
	}
	newUser := users.User{
		ID:       lastUserID + 1,
		Name:     name,
		Email:    email,
		Phone:    phone,
		IsActive: isActive,
		Scopes:   scopes,
	}
	r.Users = append(r.Users, mapToMemoryUser(newUser, password))
	return newUser, nil
}

func (r *MemoryUserRepository) Update(
	ID users.UserID, name string, email string, phone string, scopes []users.ScopeName,
) (users.User, error) {
	for i, u := range r.Users {
		if u.ID == ID {
			r.Users[i].Name = name
			r.Users[i].Email = email
			r.Users[i].Phone = phone
			r.Users[i].Scopes = scopes
			return mapToUser(r.Users[i]), nil
		}
	}
	return users.User{}, ports.UserDoesNotExists
}

func (r *MemoryUserRepository) ChangePassword(ID users.UserID, newPassword string) error {
	for i, u := range r.Users {
		if u.ID == ID {
			r.Users[i].Password = newPassword
			return nil
		}
	}
	return ports.UserDoesNotExists
}

func (r *MemoryUserRepository) Deactivate(ID users.UserID) error {
	for i, u := range r.Users {
		if u.ID == ID {
			r.Users[i].IsActive = false
			return nil
		}
	}
	return ports.UserDoesNotExists
}

func (r *MemoryUserRepository) Activate(ID users.UserID) bool {
	for i, u := range r.Users {
		if u.ID == ID {
			r.Users[i].IsActive = true
			return true
		}
	}
	return false
}

func (r *MemoryUserRepository) SaveAccountVerificationCode(ID users.UserID, code string) error {
	found := false
	for i, u := range r.Users {
		if u.ID == ID {
			r.Users[i].VerificationCode = code
			found = true
			break
		}
	}

	if !found {
		return ports.UserDoesNotExists
	}

	return nil
}

func (r *MemoryUserRepository) ValidateAccountVerificationCode(ID users.UserID, code string) bool {
	for _, u := range r.Users {
		if u.ID == ID {
			return u.VerificationCode == code
		}
	}
	return false
}

func (r *MemoryUserRepository) SaveRecoveryPasswordCode(ID users.UserID, code string) error {
	found := false
	for i, u := range r.Users {
		if u.ID == ID {
			r.Users[i].RecoveryPasswordCode = code
			found = true
			break
		}
	}

	if !found {
		return ports.UserDoesNotExists
	}

	return nil
}

func (r *MemoryUserRepository) ValidateRecoveryPasswordCode(ID users.UserID, code string) bool {
	for _, u := range r.Users {
		if u.ID == ID {
			return u.RecoveryPasswordCode == code
		}
	}
	return false
}
