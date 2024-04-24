package sqlstuff

import (
	"lms/types"

	"github.com/jmoiron/sqlx"
)

type SQLMethods interface {
	ExportData(s *sqlx.DB) error
	ImportData(s *sqlx.DB) error
	DBConnect(s *sqlx.DB) error
}

func ExportData(s sqlx.DB, users []types.User, expressions []types.Expression) error {
	_, err := s.NamedExec(`INSERT INTO users (username, encPassword, IP) VALUES (:username, :encPassword, :IP)`,
		users)
	if err != nil {
		return err
	}
	_, err = s.NamedExec(`INSERT INTO expressions (username, expression, worker, startTime, endTime) 
	VALUES (:username, :expression, :worker, :startTime, :endTime)`, expressions)
	return err
}

func ImportData(s sqlx.DB, users []types.User, expressions []types.Expression) error {
	err := s.Get(&users, `SELECT * FROM users`)
	if err != nil {
		return err
	}
	return s.Get(&expressions, `SELECT * FROM expressions`)
}
