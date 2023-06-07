package handler

import users "github.com/CrissAlvarezH/fundart-api/internal/users/domain"

type ListUserDTO struct {
	ID     users.UserID      `json:"id"`
	Name   string            `json:"name"`
	Email  string            `json:"email"`
	Phone  string            `json:"phone"`
	Scopes []users.ScopeName `json:"scopes"`
}

type RegisterUserDTO struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone"`
	Password string `json:"password" binding:"required,gte=5"`
}
