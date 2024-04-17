package httpstuff

import (
	"log"
	"slices"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Name           string
	HashedPassword string
	IP             string
}

func EncryptPassword(password, key string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"pwd": password,
	})
	tokenStr, err := token.SignedString([]byte(key))
	if err != nil {
		log.Printf("err occured: %v", err)
		panic(err)
	}
	return tokenStr
}

// Codes:
//
// 1: Everything is fine,new user signed up
//
// 0: Different user trying to sign up from same IP
//
// 2: Name already taken
func (s *Server) SignUp(password, username, key, IPAddress string) uint {
	if slices.ContainsFunc(s.Users, func(user User) bool {
		return (user.IP == IPAddress) && (user.IP != "") //empty means that person signed out
	}) {
		return 0
	} else if slices.ContainsFunc(s.Users, func(user User) bool {
		return user.Name == username
	}) {
		return 2
	} else {
		s.Users = append(s.Users, User{HashedPassword: EncryptPassword(password, key), Name: username, IP: IPAddress})
		return 1
	}
}

// Codes:
//
// 1: Everything is fine
//
// 0: Different user trying to sign in from same IP
//
// 2: Incorrect password/key
func (s *Server) SignIn(password, username, key, IPAddress string) uint {
	if slices.ContainsFunc(s.Users, func(user User) bool {
		return (user.IP == IPAddress) && (user.IP != "")
	}) {
		return 0
	} else if slices.ContainsFunc(s.Users, func(user User) bool {
		return (user.Name == username) && (user.HashedPassword == EncryptPassword(password, key))
	}) {
		s.Users[slices.IndexFunc(s.Users, func(user User) bool {
			return user.Name == username
		})].IP = IPAddress
		return 1
	} else {
		return 2
	}
}

// Codes:
//
// 1: Successfully logged out
//
// 0: No users on such IP found
func (s *Server) SignOut(IPAddress string) uint {
	ind := slices.IndexFunc(s.Users, func(user User) bool {
		return user.IP == IPAddress
	})
	if ind == -1 {
		return 0
	} else {
		s.Users[ind].IP = ""
		return 1
	}
}
