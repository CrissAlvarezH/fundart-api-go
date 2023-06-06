package handler

import users "github.com/CrissAlvarezH/fundart-api/internal/users/domain"

type ListUserDTO struct {
	ID     users.UserID      `json:"id"`
	Name   string            `json:"name"`
	Email  string            `json:"email"`
	Phone  string            `json:"phone"`
	Scopes []users.ScopeName `json:"scopes"`
}
