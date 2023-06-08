package handler

import users "github.com/CrissAlvarezH/fundart-api/internal/users/domain"

func MapToListUserDTO(user users.User) ListUserDTO {
	return ListUserDTO{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Phone:  user.Phone,
		Scopes: user.Scopes,
	}
}

func MapToListAddressDTO(address users.Address) ListAddressDTO {
	return ListAddressDTO{
		ID:            address.ID,
		Department:    address.Department,
		City:          address.City,
		Address:       address.Address,
		ReceiverPhone: address.ReceiverPhone,
		ReceiverName:  address.ReceiverName,
	}
}

func MapToListAddressesDTO(address []users.Address) []ListAddressDTO {
	dtos := make([]ListAddressDTO, 0, len(address))
	for _, a := range address {
		dtos = append(dtos, ListAddressDTO{
			ID:            a.ID,
			Department:    a.Department,
			City:          a.City,
			Address:       a.Address,
			ReceiverPhone: a.ReceiverPhone,
			ReceiverName:  a.ReceiverName,
		})
	}
	return dtos
}
