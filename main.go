package main

import (
	"lms/server"
)

func main() {
	server := server.MainServer{}
	server.Start()
}
