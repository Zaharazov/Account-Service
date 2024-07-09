package mongodb_s

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

func SaveStudent(user_id uuid.UUID, photo, full_name, group, record_book, description, mail, github string) (uuid.UUID, error) {

	student := models.Student{
		UserId:      user_id,
		Photo:       photo,
		FullName:    full_name,
		Group:       group,
		RecordBook:  record_book,
		Description: description,
		Mail:        mail,
		GitHub:      github,
	}

	// логика сохранения юзера в бд
	insertResult, err := mongodb.StudentCollection.InsertOne(context.TODO(), student)
	log.Println(insertResult)
	if err != nil {
		return user_id, err
	}

	return student.UserId, nil
}

func GetStudents(id uuid.UUID, limit int64) ([]models.Student, error) {

	options := options.Find()
	options.SetLimit(limit)

	filter := bson.D{{"userid", id}}

	if limit > 1 {
		filter = bson.D{}
	}

	var students []models.Student

	cur, err := mongodb.StudentCollection.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem models.Student
		err := cur.Decode(&elem)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		students = append(students, elem)
	}

	return students, nil
}

func GetStudentsByParams(full_name, group, record_book, description, mail, github string, limit int64) ([]models.Student, error) {

	options := options.Find()
	options.SetLimit(limit)

	onefilter := bson.D{}
	if len(full_name) > 0 {
		onefilter = append(onefilter, bson.E{"fullname", full_name})
	}
	if len(group) > 0 {
		onefilter = append(onefilter, bson.E{"group", group})
	}
	if len(record_book) > 0 {
		onefilter = append(onefilter, bson.E{"recordbook", record_book})
	}
	if len(description) > 0 {
		onefilter = append(onefilter, bson.E{"description", description})
	}
	if len(mail) > 0 {
		onefilter = append(onefilter, bson.E{"mail", mail})
	}
	if len(github) > 0 {
		onefilter = append(onefilter, bson.E{"github", github})
	}

	filter := onefilter

	var students []models.Student

	cur, err := mongodb.StudentCollection.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem models.Student
		err := cur.Decode(&elem)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		students = append(students, elem)
	}

	return students, nil
}

func DeleteStudent(id uuid.UUID) error {

	filter := bson.D{{"userid", id}}
	result, err := mongodb.StudentCollection.DeleteOne(context.TODO(), filter)

	if err != nil {
		log.Println(err)
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("Student Not Found")
	}

	return nil
}

func EditStudent(user_id uuid.UUID, photo, full_name, group, record_book, description, mail, github string) (uuid.UUID, error) {
	filter := bson.D{{"userid", user_id}}

	update := bson.D{
		{"$set", bson.D{
			{"photo", photo},
			{"fullname", full_name},
			{"group", group},
			{"recordbook", record_book},
			{"description", description},
			{"mail", mail},
			{"github", github},
		}},
	}

	result, err := mongodb.StudentCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println(err)
		return user_id, err
	}

	if result.ModifiedCount == 0 {
		return user_id, errors.New("No Changes Have Been Made")
	}

	return user_id, nil
}
