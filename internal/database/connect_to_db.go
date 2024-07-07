package database

import (
	"Account-Service/internal/configs"
	"Account-Service/internal/database/mongodb"
	"Account-Service/internal/database/postgres"
)

func ConnectToDB() {
	DBType := configs.DataBaseType
	switch DBType {
	case "mongodb":
		mongodb.ConnectToMongo()
	case "postgres":
		postgres.ConnectToPostgres()
	}
}
