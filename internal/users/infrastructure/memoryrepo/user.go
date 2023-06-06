package memoryrepo

import users "github.com/CrissAlvarezH/fundart-api/internal/users/domain"

type MemoryUserRepository struct {
	users []users.User
}

func NewMemoryUserRepository(users []users.User) *MemoryUserRepository {
	return &MemoryUserRepository{users: users}
}

func (r *MemoryUserRepository) List(
	filters map[string]string, limit int, offset int,
) ([]users.User, int, error) {
	// TODO slick users array with 'limit' and 'offset' to create pagination
	data := make([]users.User, 0)
	if len(r.users) > limit && len(r.users) > offset {
		data = r.users[offset:limit]
	} else if len(r.users) > offset && len(r.users) < limit {
		data = r.users[offset:]
	}
	return data, len(r.users), nil
}

func (r *MemoryUserRepository) GetByID(ID users.UserID) (users.User, error) {
	return users.User{}, nil
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
