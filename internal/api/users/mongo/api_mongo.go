package mongo

import (
	"context"
	"errors"
	"log"
	"restapi/internal/database/mongodb"
	"restapi/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SaveUser(login string, password string, roles []string) int32 { // TODO генерация id

	user := models.User{
		UserId:   12,
		Login:    login,
		Password: password,
		Roles:    roles,
	}

	// логика сохранения вакансии в бд
	insertResult, err := mongodb.UserCollection.InsertOne(context.TODO(), user)
	log.Println(insertResult)
	if err != nil {
		log.Fatal(err)
	}

	return user.UserId
}

func GetUsers(id int) ([]models.User, error) {

	options := options.Find()
	options.SetLimit(1)

	filter := bson.D{{"userid", id}}

	if id == -1 {
		filter = bson.D{}
	}

	var users []models.User

	cur, err := mongodb.UserCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem models.User
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		users = append(users, elem)
	}

	return users, nil
}

func DeleteUser(id int) error {

	filter := bson.D{{"userid", id}}
	result, err := mongodb.UserCollection.DeleteOne(context.TODO(), filter)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("User Not Found")
	}

	return nil
}

func EditUser(user_id int, login, password string, roles []string) (int32, error) {
	filter := bson.D{{"userid", user_id}}

	update := bson.D{
		{"$set", bson.D{
			{"login", login},
			{"password", password},
			{"roles", roles},
		}},
	}

	result, err := mongodb.UserCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	if result.ModifiedCount == 0 {
		return -1, errors.New("No Changes Have Been Made")
	}

	return int32(user_id), nil
}
