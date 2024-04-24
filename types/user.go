package types

type User struct {
	Name           string `db:"username"`
	HashedPassword string `db:"encPassword"`
	IP             string `db:"IP"`
}
