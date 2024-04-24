package httpstuff

import (
	"lms/types"
	"slices"
	"time"
)

// Codes:
//
// 1: Everything is fine,new user signed up
//
// 0: Different user trying to sign up from same IP
//
// 2: Name already taken
func (s *Server) SignUp(password, username, key, IPAddress string) uint {
	if slices.ContainsFunc(s.Users, func(user types.User) bool {
		return (user.IP == IPAddress) && (user.IP != "") //empty means that person signed out
	}) {
		return 0
	} else if slices.ContainsFunc(s.Users, func(user types.User) bool {
		return user.Name == username
	}) {
		return 2
	} else {
		s.Users = append(s.Users, types.User{HashedPassword: EncryptPassword(password, key), Name: username, IP: IPAddress})
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
	if slices.ContainsFunc(s.Users, func(user types.User) bool {
		return (user.IP == IPAddress) && (user.IP != "")
	}) {
		return 0
	} else if slices.ContainsFunc(s.Users, func(user types.User) bool {
		return (user.Name == username) && (user.HashedPassword == EncryptPassword(password, key))
	}) {
		s.Users[slices.IndexFunc(s.Users, func(user types.User) bool {
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
	ind := slices.IndexFunc(s.Users, func(user types.User) bool {
		return user.IP == IPAddress
	})
	if ind == -1 {
		return 0
	} else {
		s.Users[ind].IP = ""
		return 1
	}
}

// Codes:
//
// 1: User verified
//
// 0: User is not signed in
func (s *Server) VerifyUser(IPAddress string) (string, uint) {
	ind := slices.IndexFunc(s.Users, func(user types.User) bool {
		return user.IP == IPAddress
	})
	if ind > -1 {
		return s.Users[ind].Name, 1
	} else {
		return "", 0
	}
}

func (s *Server) GetExpressions(username string) []types.Expression {
	var result []types.Expression
	for _, val := range s.Expressions {
		if val.Username == username {
			result = append(result, val)
		}
	}
	return result
}

func (s *Server) AddExpression(expression, username string) {
	var result types.Expression
	result.StartTime = time.Now()
	result.EndTime = result.StartTime.Add(s.WorkingTime)
	result.Expression = expression
	result.Username = username

	min := ^uint(0) //the biggest uint value possible
	var workerName string
	for key, val := range s.Workers {
		if val < min {
			min = val
			workerName = key //finding the worker with the least amount of expressions
		}
	}
	result.WorkerName = workerName
}

// This function checks if it's needed to remove expression due to its timeout or not.
func (s *Server) TryRemovingExpression(expression types.Expression) {
	if time.Now().After(expression.EndTime) {
		ind := slices.Index(s.Expressions, expression)
		s.Expressions[ind] = s.Expressions[len(s.Expressions)-1]
		s.Expressions = s.Expressions[:len(s.Expressions)-1]
		s.Workers[expression.WorkerName]--
	}
}
