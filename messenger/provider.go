package messenger

import (
	"github.com/Koltsov1998/MiniPM/user"
)

type IMessengerProvider interface {
	SendMessage(userId user.Id, message string) (chan string, error)
	SendMessageWithoutResponse(userId user.Id, message string) error
}
