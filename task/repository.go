package task

import (
	"github.com/Koltsov1998/MiniPM/user"
)

type ITaskRepository[U user.User, T Task] interface {
	GetAllInProgress(users []U) (map[user.Id][]T, error)
	GetAllInProgressForUser(user U) ([]T, error)
	WriteTaskReport(task Task, user U, reportMessage string) error
}
