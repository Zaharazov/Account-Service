package mongodb_u

import (
	"Account-Service/internal/domain/models"
	mongodb "Account-Service/internal/services/database"
	mongodb_e "Account-Service/internal/services/database/employers"
	mongodb_o "Account-Service/internal/services/database/organizers"
	mongodb_s "Account-Service/internal/services/database/students"
	"context"
	"errors"
	"log"
	"sort"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var const_roles []string = []string{"student", "organizer", "employer"}

func CompareRoles(old_roles, new_roles []string) map[string]string {
	var results map[string]string
	var orMap map[string]bool
	var nrMap map[string]bool

	for _, role := range const_roles {
		orMap[role] = false
		nrMap[role] = false
		results[role] = ""
	}

	for _, role := range const_roles {
		orMap[role] = (sort.SearchStrings(old_roles, role) != len(old_roles))
		nrMap[role] = (sort.SearchStrings(old_roles, role) != len(new_roles))
	}

	for _, role := range const_roles {
		or := orMap[role]
		nr := nrMap[role]
		if or == nr {
			results[role] = "nothing" // or = nr = 1 / or = nr = 0
		} else if !or && nr {
			results[role] = "create" // or = 0, nr = 1
		} else if or && !nr {
			results[role] = "delete" // or = 1, nr = 0
		}
	}

	return results
}

func SaveUser(user_id uuid.UUID, login string, password string, roles []string) (uuid.UUID, error) { // TODO генерация id (используем пакет uuid от гугла)

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

	new_roles := roles

	users, err := GetUsers(user_id, 1)
	old_roles := users[0].Roles

	results := CompareRoles(old_roles, new_roles)

	result, err := mongodb.UserCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println(err)
		return user_id, err
	}

	s_status := results["student"]
	o_status := results["organizer"]
	e_status := results["employer"]

	switch s_status {
	case "create":
		_, err = mongodb_s.SaveStudent(user_id, "", "", "", "", "", "", "")
		if err != nil {
			return user_id, err
		}
	case "delete":
		err = mongodb_s.DeleteStudent(user_id)
		if err != nil {
			return user_id, err
		}
	}

	switch o_status {
	case "create":
		_, err = mongodb_o.SaveOrganizer(user_id, "", "", "")
		if err != nil {
			return user_id, err
		}
	case "delete":
		err = mongodb_o.DeleteOrganizer(user_id)
		if err != nil {
			return user_id, err
		}
	}

	switch e_status {
	case "create":
		_, err = mongodb_e.SaveEmployer(user_id, "", "", "", "")
		if err != nil {
			return user_id, err
		}
	case "delete":
		err = mongodb_e.DeleteEmployer(user_id)
		if err != nil {
			return user_id, err
		}
	}

	if result.ModifiedCount == 0 {
		return user_id, errors.New("No Changes Have Been Made")
	}

	return user_id, nil
}
