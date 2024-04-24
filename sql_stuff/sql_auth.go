package sqlstuff

import (
	"os"
	"path"
	"strings"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func DBConnect() (*sqlx.DB, error) {
	curr_dir, _ := os.Getwd() //the function is called only in lms/cmd/main.go
	par_dir := path.Dir(curr_dir)
	data, err := os.ReadFile(par_dir + "/DBInfo.txt")
	if err != nil {
		return nil, err
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
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	return db, err
}
