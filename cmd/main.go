package main

import (
	"log"
	"net/http"
	"restapi/internal/configs"
	"restapi/internal/database/postgres"
	"restapi/internal/routers"
)

func main() {

	log.Printf("Server started")

	postgres.ConnectToPostgres()
	//mongodb.ConnectToMongo()

	router := routers.NewRouter()

	log.Fatal(http.ListenAndServe(configs.Port, router))
}
