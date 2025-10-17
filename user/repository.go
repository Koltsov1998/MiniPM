package user

import "time"

type Id int64

// Минимальный контракт, общий для всех "пользователей".
type User interface {
	GetId() Id
	GetName() string
	GetLastTimeNotified() time.Time
}

// Дженерик-репозиторий под любой U, удовлетворяющий User.
type IUserRepository[U User] interface {
	GetAll() ([]U, error)
	GetById(id Id) (U, error)
	Create(u U) (U, error)
}
