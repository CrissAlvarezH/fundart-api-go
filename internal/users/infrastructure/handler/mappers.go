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
