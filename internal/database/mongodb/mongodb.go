package mongodb

import (
	"Account-Service/internal/configs"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client, err = mongo.NewClient(options.Client().ApplyURI(configs.MongoURI))

var UserCollection = client.Database(configs.DBName).Collection(configs.UsersCollectionName)
var StudentCollection = client.Database(configs.DBName).Collection(configs.StudentsCollectionName)
var EmployerCollection = client.Database(configs.DBName).Collection(configs.EmployersCollectionName)
var OrganizerCollection = client.Database(configs.DBName).Collection(configs.OrganizersCollectionName)

func ConnectToMongo() {
	// Create connect
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
}
