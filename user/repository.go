package user

import "time"

type Id int64

type User interface {
	GetId() Id
	GetName() string
	GetLastTimeNotified() time.Time
}

type IUserRepository[U User] interface {
	GetAll() ([]U, error)
	GetById(id Id) (U, error)
	Create(u U) (U, error)
}
