package httpstuff

import (
	"encoding/json"
	sql "lms/sql_stuff"
	"lms/types"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

type Server struct {
	Mu          *sync.Mutex
	WorkingTime time.Duration
	Workers     map[string]uint
	Mux         *http.ServeMux
	Users       []types.User
	DBConn      *sqlx.DB
	Expressions []types.Expression
}

type WebMethods interface {
	Start(firstTime bool) error
	LogIn(password, username, key, IPAddress string) uint
	SignUp(password, username, key, IPAddress string) uint
	SignOut(IPAddress string) uint
	VerifyUser(IPAddress, username string) uint
}

func (s *Server) GetWorkersNames() error {
	curr_dir, _ := os.Getwd() //the function is called only in lms/cmd/main.go
	par_dir := path.Dir(curr_dir)
	data, err := os.ReadFile(par_dir + "/DBInfo.txt")
	if err != nil {
		return err
	}
	names := strings.Split(string(data), "\n")
	result := make(map[string]uint, len(names))
	for _, name := range names {
		result[name] = 0
	}
	s.Workers = result
	return nil
}

func (s *Server) Start(firstTime bool) error {

	if err := s.GetWorkersNames(); err != nil {
		return err
	}

	if !firstTime {
		err := sql.ImportData(*s.DBConn, s.Users, s.Expressions)
		if err != nil {
			return err
		}
		for _, exp := range s.Expressions {
			s.Workers[exp.WorkerName]++
		}
	}

	go func() {
		for range time.Tick(time.Second * 10) { //autosave data every 10 seconds
			s.Mu.Lock()
			sql.ExportData(*s.DBConn, s.Users, s.Expressions)
			s.Mu.Unlock()
		}
	}()

	go func() {
		for range time.Tick(time.Second * 3) { //remove "expired" expressions every 3 seconds
			s.Mu.Lock()
			for _, val := range s.Expressions {
				s.TryRemovingExpression(val)
			}
			s.Mu.Unlock()
		}
	}()

	s.Mux.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		var (
			IP       = r.RemoteAddr
			Password = r.URL.Query().Get("password")
			Key      = r.URL.Query().Get("key")
			Username = r.URL.Query().Get("name")
		)
		code := s.SignUp(Password, Username, Key, IP)
		switch code {
		case 1:
			w.Write([]byte("Everything fine!"))
		case 2:
			w.Write([]byte("This username is already taken"))
		default:
			w.Write([]byte("You've already signed up/in from this device!"))
		}
	})

	s.Mux.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		var (
			IP       = r.RemoteAddr
			Password = r.URL.Query().Get("password")
			Key      = r.URL.Query().Get("key")
			Username = r.URL.Query().Get("name")
		)
		code := s.SignIn(Password, Username, Key, IP)
		switch code {
		case 1:
			w.Write([]byte("Everything fine!"))
		case 0:
			w.Write([]byte("Incorrect password and key combination, try again!"))
		default:
			w.Write([]byte("You've already signed up/in from this device!"))
		}
	})

	s.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		name, code := s.VerifyUser(r.RemoteAddr)
		if code != 1 {
			w.Write([]byte("You're not signed in!"))
			return
		} else {
			reply, err := json.Marshal(s.GetExpressions(name))
			if err != nil {
				w.Write([]byte(err.Error()))
			} else {
				w.Write(reply)
			}
		}
	})

	s.Mux.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		name, code := s.VerifyUser(r.RemoteAddr)
		if code != 1 {
			w.Write([]byte("You're not signed in!"))
			return
		} else {
			s.AddExpression(r.URL.Query().Get("exp"), name)
			w.Write([]byte("new expression added!"))
		}
	})

	s.Mux.HandleFunc("/signout", func(w http.ResponseWriter, r *http.Request) {
		_, code := s.VerifyUser(r.RemoteAddr)
		if code != 1 {
			w.Write([]byte("You're not signed in!"))
			return
		} else {
			s.SignOut(r.RemoteAddr)
			w.Write([]byte("Successfully logged out!"))
		}
	})

	return http.ListenAndServe(":8000", s.Mux)
}
