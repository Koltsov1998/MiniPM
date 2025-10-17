package messenger

import (
	"github.com/Koltsov1998/MiniPM/user"
)

type IMessengerProvider[U user.User] interface {
	SendMessage(user U, message string) (chan string, error)
	SendMessageWithoutResponse(user U, message string) error
}
