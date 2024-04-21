package httpstuff

import (
	"database/sql"
)

type Server struct {
	Users  []User
	DBConn *sql.DB
}

type WebMethods interface {
	LogIn(password, username, key, IPAddress string) uint
	SignUp(password, username, key, IPAddress string) uint
	SignOut(IPAddress string) uint
	VerifyUser(IPAddress, username string) uint
}
