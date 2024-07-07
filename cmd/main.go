package main

import (
	"log"
	"net/http"
	"restapi/internal/configs"
	"restapi/internal/database/mongodb"
	"restapi/internal/routers"
)

func main() {

	log.Printf("Server started")

	mongodb.ConnectToMongo()

	router := routers.NewRouter()

	log.Fatal(http.ListenAndServe(configs.HttpPort, router))
}
