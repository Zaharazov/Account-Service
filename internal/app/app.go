package app

import (
	"Account-Service/internal/configs"
	"Account-Service/internal/database"
	"Account-Service/internal/routers"
	"log"
	"net/http"
)

func Run() {
	log.Printf("Server started")

	database.ConnectToDB()

	router := routers.NewRouter()

	log.Fatal(http.ListenAndServe(configs.Port, router))
}
