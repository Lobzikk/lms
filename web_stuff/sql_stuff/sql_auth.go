package sqlstuff

import (
	"database/sql"
	web "lms/web_stuff"
	"os"
	"path"
	"strings"

	_ "github.com/lib/pq"
)

func DBConnect(server *web.Server) error {
	curr_dir, _ := os.Getwd() //the function is called only in lms/cmd/main.go
	par_dir := path.Dir(curr_dir)
	data, err := os.ReadFile(par_dir + "/DBInfo.txt")
	if err != nil {
		return err
	}
	DBAuthData := strings.Split(string(data), "\n")
	var ( // look for lms/setup.bat
		host     = strings.Split(DBAuthData[0], ":")[0]
		port     = strings.Split(DBAuthData[0], ":")[1]
		user     = DBAuthData[1]
		password = DBAuthData[2]
		dbname   = DBAuthData[3]
	)
	psqlInfo := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	server.DBConn = db
	err = server.DBConn.Ping()
	return err
}
