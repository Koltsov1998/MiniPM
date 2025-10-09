package user

type UserId int64

type User struct {
	ID              UserId `json:"id"`
	Name            string `json:"name"`
	ReddyID         string `json:"reddy_id"`
	LastDayNotified string `json:"last_day_notified"`
}

type IUserRepository interface {
	GetAll() ([]User, error)
}
