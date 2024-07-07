package mongodb

import (
	"Account-Service/internal/database/mongodb"
	"Account-Service/internal/models"
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SaveUser(login string, password string, roles []string) (uuid.UUID, error) { // TODO генерация id (используем пакет uuid от гугла)

	user_id := uuid.New()

	user := models.User{
		UserId:   user_id,
		Login:    login,
		Password: password,
		Roles:    roles,
	}

	// логика сохранения юзера в бд
	insertResult, err := mongodb.UserCollection.InsertOne(context.TODO(), user)
	log.Println(insertResult)
	if err != nil {
		return user_id, err
	}

	return user.UserId, nil
}

func GetUsers(id uuid.UUID, limit int64) ([]models.User, error) {

	options := options.Find()
	options.SetLimit(limit)

	filter := bson.D{{"userid", id}}

	if limit > 1 {
		filter = bson.D{}
	}

	var users []models.User

	cur, err := mongodb.UserCollection.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem models.User
		err := cur.Decode(&elem)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		users = append(users, elem)
	}

	return users, nil
}

func DeleteUser(id uuid.UUID) error {

	filter := bson.D{{"userid", id}}
	result, err := mongodb.UserCollection.DeleteOne(context.TODO(), filter)

	if err != nil {
		log.Println(err)
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("User Not Found")
	}

	return nil
}

func EditUser(user_id uuid.UUID, login, password string, roles []string) (uuid.UUID, error) {
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
		log.Println(err)
		return user_id, err
	}

	if result.ModifiedCount == 0 {
		return user_id, errors.New("No Changes Have Been Made")
	}

	return user_id, nil
}
