package memoryrepo

import (
	"github.com/CrissAlvarezH/fundart-api/internal/users/application/ports"
	users "github.com/CrissAlvarezH/fundart-api/internal/users/domain"
)

type MemoryAddress struct {
	ID            users.AddressID
	Department    string
	City          string
	Address       string
	ReceiverPhone string
	ReceiverName  string
	UserID        users.UserID
}

func mapToMemoryAddress(address users.Address, userID users.UserID) MemoryAddress {
	return MemoryAddress{
		ID:            address.ID,
		Department:    address.Department,
		City:          address.City,
		Address:       address.Address,
		ReceiverPhone: address.ReceiverPhone,
		ReceiverName:  address.ReceiverName,
		UserID:        userID,
	}
}

func mapToAddress(memoAddress MemoryAddress) users.Address {
	return users.Address{
		ID:            memoAddress.ID,
		Department:    memoAddress.Department,
		City:          memoAddress.City,
		Address:       memoAddress.Address,
		ReceiverPhone: memoAddress.ReceiverPhone,
		ReceiverName:  memoAddress.ReceiverName,
	}
}

type MemoryAddressRepository struct {
	addresses []MemoryAddress
}

func NewMemoryAddressRepository(addresses []MemoryAddress) *MemoryAddressRepository {
	return &MemoryAddressRepository{addresses: addresses}
}

func (r *MemoryAddressRepository) List(userID users.UserID) []users.Address {
	var addresses = make([]users.Address, 0)
	for _, a := range r.addresses {
		if a.UserID == userID {
			addresses = append(addresses, mapToAddress(a))
		}
	}
	return addresses
}

func (r *MemoryAddressRepository) Add(
	userID users.UserID, department string, city string, address string, receiverPhone string, receiverName string,
) (users.Address, error) {
	lastAddressID := 0
	if len(r.addresses) > 0 {
		lastAddressID = int(r.addresses[len(r.addresses)-1].ID)
	}

	newAddress := users.Address{
		ID:            users.AddressID(lastAddressID + 1),
		Department:    department,
		City:          city,
		Address:       address,
		ReceiverPhone: receiverPhone,
		ReceiverName:  receiverName,
	}
	r.addresses = append(r.addresses, mapToMemoryAddress(newAddress, userID))
	return newAddress, nil
}

func (r *MemoryAddressRepository) Update(
	ID users.AddressID, department string, city string, address string,
	receiverPhone string, receiverName string,
) (users.Address, error) {
	for i, u := range r.addresses {
		if u.ID == ID {
			r.addresses[i].Department = department
			r.addresses[i].City = city
			r.addresses[i].Address = address
			r.addresses[i].ReceiverName = receiverPhone
			r.addresses[i].ReceiverPhone = receiverName
			return mapToAddress(r.addresses[i]), nil
		}
	}
	return users.Address{}, ports.AddressDoesNotExists
}

func (r *MemoryAddressRepository) Delete(ID users.AddressID) error {
	found := false
	filtered := make([]MemoryAddress, 0, len(r.addresses))
	for _, a := range r.addresses {
		if a.ID != ID {
			filtered = append(filtered, a)
			found = true
		}
	}

	r.addresses = filtered
	if found == false {
		return ports.AddressDoesNotExists
	}
	return nil
}
