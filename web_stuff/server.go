package httpstuff

type Server struct {
	Users []User
}

type WebMethods interface {
	LogIn(password, username, key, IPAddress string) uint
	SignUp(password, username, key, IPAddress string) uint
	SignOut(IPAddress string) uint
}
