package main

import (
	"Account-Service/internal/configs"
	"Account-Service/internal/database"
	"Account-Service/internal/routers"
	"log"
	"net/http"
)

func main() {

	log.Printf("Server started")

	database.ConnectToDB()

	router := routers.NewRouter()

	log.Fatal(http.ListenAndServe(configs.Port, router))
}
