package main

import (
	"log"
	"net/http"
	"restapi/integer/database/mongodb"
	"restapi/integer/routers"
)

func main() {

	log.Printf("Server started")

	mongodb.ConnectToMongo()

	router := routers.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
