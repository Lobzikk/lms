package main

import (
	web "lms/web_stuff"
	"log"
)

func main() {
	var server web.Server
	if err := server.Start(true); err != nil {
		panic(err)
	} else {
		log.Print("Server started listening on port 8000!")
	}
}
