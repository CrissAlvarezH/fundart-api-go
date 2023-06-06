package users

import "time"

type UserID int
type ScopeName string

type User struct {
	ID        UserID
	Name      string
	Email     string
	Phone     string
	IsActive  bool
	CreatedAt time.Time
	Addresses []Address
	Scopes    []ScopeName
}

type AddressID int

type Address struct {
	ID            AddressID
	Department    string
	City          string
	Address       string
	ReceiverPhone string
	ReceiverName  string
}
