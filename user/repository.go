package user

import "time"

type Id int64

type User struct {
	ID               Id        `json:"id"`
	Name             string    `json:"name"`
	LastTimeNotified time.Time `json:"last_time_notified"`
}

type IUserRepository[U User] interface {
	GetAll() ([]*U, error)
	GetById(id Id) (*U, error)
	Create(User) (*U, error)
}
