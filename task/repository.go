package task

import (
	"MiniPm/user"
)

type Task struct {
	Title string
}

type ITaskRepository interface {
	GetAllInProgress(userIds []user.UserId) (map[user.UserId][]Task, error)
	GetAllInProgressForUser(userId user.UserId) ([]Task, error)
	WriteTaskReport(userId user.UserId, reportMessage string) error
}
