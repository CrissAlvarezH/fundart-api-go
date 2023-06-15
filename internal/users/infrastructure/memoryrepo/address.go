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
	Addresses []MemoryAddress
}

func NewMemoryAddressRepository(addresses []MemoryAddress) *MemoryAddressRepository {
	return &MemoryAddressRepository{Addresses: addresses}
}

func (r *MemoryAddressRepository) List(userID users.UserID) []users.Address {
	var addresses = make([]users.Address, 0)
	for _, a := range r.Addresses {
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
	if len(r.Addresses) > 0 {
		lastAddressID = int(r.Addresses[len(r.Addresses)-1].ID)
	}

	newAddress := users.Address{
		ID:            users.AddressID(lastAddressID + 1),
		Department:    department,
		City:          city,
		Address:       address,
		ReceiverPhone: receiverPhone,
		ReceiverName:  receiverName,
	}
	r.Addresses = append(r.Addresses, mapToMemoryAddress(newAddress, userID))
	return newAddress, nil
}

func (r *MemoryAddressRepository) Update(
	ID users.AddressID, department string, city string, address string,
	receiverPhone string, receiverName string,
) (users.Address, error) {
	for i, u := range r.Addresses {
		if u.ID == ID {
			r.Addresses[i].Department = department
			r.Addresses[i].City = city
			r.Addresses[i].Address = address
			r.Addresses[i].ReceiverName = receiverPhone
			r.Addresses[i].ReceiverPhone = receiverName
			return mapToAddress(r.Addresses[i]), nil
		}
	}
	return users.Address{}, ports.AddressDoesNotExists
}

func (r *MemoryAddressRepository) Delete(ID users.AddressID) error {
	found := false
	filtered := make([]MemoryAddress, 0, len(r.Addresses))
	for _, a := range r.Addresses {
		if a.ID != ID {
			filtered = append(filtered, a)
			found = true
		}
	}

	r.Addresses = filtered
	if !found {
		return ports.AddressDoesNotExists
	}
	return nil
}
