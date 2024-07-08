package database

import (
	"Account-Service/internal/database/mongodb"
)

func ConnectToDB() {
	mongodb.ConnectToMongo()
}
