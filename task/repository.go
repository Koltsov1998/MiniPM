package task

import (
	"github.com/Koltsov1998/MiniPM/user"
)

type ITaskRepository interface {
	GetAllInProgress(users []user.User) (map[user.UserId][]Task, error)
	GetAllInProgressForUser(user user.User) ([]Task, error)
	WriteTaskReport(user user.User, reportMessage string) error
}
