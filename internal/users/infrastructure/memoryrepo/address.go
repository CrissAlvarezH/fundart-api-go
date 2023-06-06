package memoryrepo

import users "github.com/CrissAlvarezH/fundart-api/internal/users/domain"

type MemoryAddressRepository struct {
	addresses []users.Address
}

func NewMemoryAddressRepository(addresses []users.Address) *MemoryAddressRepository {
	return &MemoryAddressRepository{addresses: addresses}
}

func (r *MemoryAddressRepository) Add(
	department string, city string, address string, receiverPhone string, receiverName string,
) (users.Address, error) {
	return users.Address{}, nil
}

func (r *MemoryAddressRepository) Update(
	ID users.AddressID, department string, city string, address string,
	receiverPhone string, receiverName string,
) (users.Address, error) {
	return users.Address{}, nil
}

func (r *MemoryAddressRepository) Delete(ID users.AddressID) error {
	return nil
}
