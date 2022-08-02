package main

import (
	"FireBaseEx/server"
	"log"
)

func main() {
	srv := server.SetupRoutes()
	log.Println("Connected")
	err := srv.Run(":8081")
	if err != nil {
		log.Fatal(err)
		return
	}

}
