package organizers

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

func SaveOrganizer(user_id uuid.UUID, photo, name, description string) (uuid.UUID, error) {

	organizer := models.Organizer{
		UserId:        user_id,
		Photo:         photo,
		Name:          name,
		Description:   description,
		CreatedEvents: nil,
	}

	// логика сохранения юзера в бд
	insertResult, err := mongodb.OrganizerCollection.InsertOne(context.TODO(), organizer)
	log.Println(insertResult)
	if err != nil {
		return user_id, err
	}

	return organizer.UserId, nil
}

func GetOrganizers(id uuid.UUID, limit int64) ([]models.Organizer, error) {

	options := options.Find()
	options.SetLimit(limit)

	filter := bson.D{{"userid", id}}

	if limit > 1 {
		filter = bson.D{}
	}

	var organizers []models.Organizer

	cur, err := mongodb.OrganizerCollection.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem models.Organizer
		err := cur.Decode(&elem)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		organizers = append(organizers, elem)
	}

	return organizers, nil
}

func GetOrganizersByParams(name, description string, limit int64) ([]models.Organizer, error) {

	options := options.Find()
	options.SetLimit(limit)

	onefilter := bson.D{}
	if len(name) > 0 {
		onefilter = append(onefilter, bson.E{"name", name})
	}
	if len(description) > 0 {
		onefilter = append(onefilter, bson.E{"description", description})
	}

	filter := onefilter

	var organizers []models.Organizer

	cur, err := mongodb.OrganizerCollection.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem models.Organizer
		err := cur.Decode(&elem)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		organizers = append(organizers, elem)
	}

	return organizers, nil
}

func DeleteOrganizer(id uuid.UUID) error {

	filter := bson.D{{"userid", id}}
	result, err := mongodb.OrganizerCollection.DeleteOne(context.TODO(), filter)

	if err != nil {
		log.Println(err)
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("Organizer Not Found")
	}

	return nil
}

func EditOrganizer(user_id uuid.UUID, photo, name, description string) (uuid.UUID, error) {
	filter := bson.D{{"userid", user_id}}

	update := bson.D{
		{"$set", bson.D{
			{"photo", photo},
			{"name", name},
			{"description", description},
		}},
	}

	result, err := mongodb.OrganizerCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println(err)
		return user_id, err
	}

	if result.ModifiedCount == 0 {
		return user_id, errors.New("No Changes Have Been Made")
	}

	return user_id, nil
}
