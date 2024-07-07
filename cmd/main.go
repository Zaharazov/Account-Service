package main

import (
	"Account-Service/internal/configs"
	"Account-Service/internal/database/mongodb"
	"Account-Service/internal/routers"
	"log"
	"net/http"
)

func main() {

	log.Printf("Server started")

	// DBType := configs.DataBaseType
	// switch DBType {
	// case "mongodb":
	// 	mongodb.ConnectToMongo()
	// case "postgres":
	// 	postgres.ConnectToPostgres()
	// }

	mongodb.ConnectToMongo()

	router := routers.NewRouter()

	log.Fatal(http.ListenAndServe(configs.Port, router))
}
