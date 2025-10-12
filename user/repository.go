package user

import "time"

type UserId int64

type User struct {
	ID               UserId    `json:"id"`
	Name             string    `json:"name"`
	LastTimeNotified time.Time `json:"last_time_notified"`
}

type IUserRepository interface {
	GetAll() ([]*User, error)
	GetById(id UserId) (*User, error)
	Create(User) (*User, error)
}
