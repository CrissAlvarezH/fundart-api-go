package handler

import users "github.com/CrissAlvarezH/fundart-api/internal/users/domain"

type ListUserDTO struct {
	ID     users.UserID      `json:"id"`
	Name   string            `json:"name"`
	Email  string            `json:"email"`
	Phone  string            `json:"phone"`
	Scopes []users.ScopeName `json:"scopes"`
}

type RetrieveUserDTO struct {
	ID        users.UserID      `json:"id"`
	Name      string            `json:"name"`
	Email     string            `json:"email"`
	Phone     string            `json:"phone"`
	Scopes    []users.ScopeName `json:"scopes"`
	Addresses []ListAddressDTO  `json:"addresses"`
}

type RegisterUserDTO struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone"`
	Password string `json:"password" binding:"required,gte=5"`
}

type LoginUserDTO struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,gte=5"`
}

type ValidateVerificationCodeDTO struct {
	Code string `json:"code" binding:"required"`
}

type RequestRecoveryPasswordDTO struct {
	Email string `json:"email" binding:"required"`
}

type RecoveryPasswordDTO struct {
	Email       string `json:"email" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,gte=5"`
	Code        string `json:"code" binding:"required"`
}

type UpdateUserDTO struct {
	Name   string            `json:"name" binding:"required"`
	Email  string            `json:"email" binding:"required"`
	Phone  string            `json:"phone"`
	Scopes []users.ScopeName `json:"scopes"`
}

type ChangePasswordDTO struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type ListAddressDTO struct {
	ID            users.AddressID `json:"id"`
	Department    string          `json:"department"`
	City          string          `json:"city"`
	Address       string          `json:"address"`
	ReceiverPhone string          `json:"receiver_phone"`
	ReceiverName  string          `json:"receiver_name"`
}

type CreateAddressDTO struct {
	Department    string `json:"department"`
	City          string `json:"city"`
	Address       string `json:"address"`
	ReceiverPhone string `json:"receiver_phone"`
	ReceiverName  string `json:"receiver_name"`
}
