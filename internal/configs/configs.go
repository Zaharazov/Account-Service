package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	HttpPort                 = GetConfigs("httpPort")
	MongoURI                 = GetConfigs("mongoURI")
	DBName                   = GetConfigs("databaseName")
	UsersCollectionName      = GetConfigs("usersCollectionName")
	StudentsCollectionName   = GetConfigs("studentsCollectionName")
	EmployersCollectionName  = GetConfigs("employersCollectionName")
	OrganizersCollectionName = GetConfigs("organizersCollectionName")
	Port                     = GetConfigs("httpPort")
)

func GetConfigs(param string) string {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Data, exists := os.LookupEnv(param)

	if exists {
		Data = os.Getenv(param)
		log.Printf("%s is %s", param, Data)
	} else {
		log.Printf("%s is missing", param)
	}

	return Data

}
