package employers

import (
	"Account-Service/internal/domain/models"
	mongodb "Account-Service/internal/services/database"
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SaveEmployer(user_id uuid.UUID, photo, name, focus, description string) (uuid.UUID, error) {

	employer := models.Employer{
		UserId:           user_id,
		Photo:            photo,
		Name:             name,
		Focus:            focus,
		Description:      description,
		CreatedVacancies: nil,
	}

	// логика сохранения юзера в бд
	insertResult, err := mongodb.EmployerCollection.InsertOne(context.TODO(), employer)
	log.Println(insertResult)
	if err != nil {
		return user_id, err
	}

	return employer.UserId, nil
}

func GetEmployers(id uuid.UUID, limit int64) ([]models.Employer, error) {

	options := options.Find()
	options.SetLimit(limit)

	filter := bson.D{{"userid", id}}

	if limit > 1 {
		filter = bson.D{}
	}

	var employers []models.Employer

	cur, err := mongodb.EmployerCollection.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem models.Employer
		err := cur.Decode(&elem)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		employers = append(employers, elem)
	}

	return employers, nil
}

func GetEmployersByParams(name, focus, description string, limit int64) ([]models.Employer, error) {

	options := options.Find()
	options.SetLimit(limit)

	onefilter := bson.D{}
	if len(name) > 0 {
		onefilter = append(onefilter, bson.E{"name", name})
	}
	if len(focus) > 0 {
		onefilter = append(onefilter, bson.E{"focus", focus})
	}
	if len(description) > 0 {
		onefilter = append(onefilter, bson.E{"description", description})
	}

	filter := onefilter

	var employers []models.Employer

	cur, err := mongodb.EmployerCollection.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem models.Employer
		err := cur.Decode(&elem)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		employers = append(employers, elem)
	}

	return employers, nil
}

func DeleteEmployer(id uuid.UUID) error {

	filter := bson.D{{"userid", id}}
	result, err := mongodb.EmployerCollection.DeleteOne(context.TODO(), filter)

	if err != nil {
		log.Println(err)
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("Employer Not Found")
	}

	return nil
}

func EditEmployer(user_id uuid.UUID, photo, focus, name, description string) (uuid.UUID, error) {
	filter := bson.D{{"userid", user_id}}

	update := bson.D{
		{"$set", bson.D{
			{"photo", photo},
			{"name", name},
			{"focus", focus},
			{"description", description},
		}},
	}

	result, err := mongodb.EmployerCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println(err)
		return user_id, err
	}

	if result.ModifiedCount == 0 {
		return user_id, errors.New("No Changes Have Been Made")
	}

	return user_id, nil
}
