package app

import (
	"Account-Service/internal/presentation/routers"
	"Account-Service/internal/server/configs"
	mongodb "Account-Service/internal/services/database"
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func Run() {
	log.Printf("Server started")

	mongodb.ConnectToMongo()

	router := routers.NewRouter()
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(httpSwagger.URL("http://localhost" + configs.Port + "/swagger/doc.json")))

	log.Fatal(http.ListenAndServe(configs.Port, router))
}
