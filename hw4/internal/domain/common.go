package domain

const (
	ServerAddr = ":5000"

	Users    = "/users"
	Login    = "/login"
	Register = "/register"

	AllMessages = "/messages"
	NumMessages = "/messages/{num}"
	PrivateMsg  = "/messages/my"
	PrivateTo   = "/messages/{to_user}"
)

type Message struct {
	User string
	Msg  string
}

type UserData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewMessage(user string, message string) Message {
	return Message{
		User: user,
		Msg:  message,
	}
}
