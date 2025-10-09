package messenger

import (
	"MiniPm/user"
)

type IMessengerProvider interface {
	SendMessage(userId user.UserId, message string) (chan string, error)
	SendMessageWithoutResponse(userId user.UserId, message string) error
}
