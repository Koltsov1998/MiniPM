package messenger

import (
	"github.com/Koltsov1998/MiniPM/user"
)

type IMessengerProvider interface {
	SendMessage(userId user.UserId, message string) (chan string, error)
	SendMessageWithoutResponse(userId user.UserId, message string) error
}
