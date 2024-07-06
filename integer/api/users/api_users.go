package users

import (
	"context"
	"encoding/json"
	"log"
	"mime"
	"net/http"
	"restapi/integer/database/mongodb"
	"restapi/integer/models"
	"strconv"

	"github.com/gorilla/mux"
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

func GetUser(id int) (models.User, error) {

	options := options.Find()
	options.SetLimit(1)

	filter := bson.D{{"userid", id}}

	var user models.User

	err := mongodb.UserCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	type RequestUser struct {
		Login    string   `json:"login,omitempty"`
		Password string   `json:"password,omitempty"`
		Roles    []string `json:"roles,omitempty"`
	}

	type ResponseUserId struct {
		Id int32 `json:"id"` // структура объекта, который отдаем
	}

	contentType := r.Header.Get("Content-Type")           // получаем тип контента в запросе
	mediatype, _, err := mime.ParseMediaType(contentType) // парсинг полученных данных из запроса
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // ошибка, если не получили запрос
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType) // ошибка, если пришел не json
		return
	}

	dec := json.NewDecoder(r.Body) // декодируем тело запроса
	dec.DisallowUnknownFields()    // проверка, что получили то, что готовы принять (ругается на id, если оно есть в запросе)
	var ru RequestUser
	if err := dec.Decode(&ru); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // получаем декодированные данные и проверяем, что все ок
		return
	}

	if len(ru.Login) == 0 || len(ru.Password) == 0 || ru.Roles == nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := SaveUser(ru.Login, ru.Password, ru.Roles) // создаем карточку по данным из запроса и получаем его id
	js, err := json.Marshal(ResponseUserId{Id: id}) // формируем json ответ с id выше
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func DeleteUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	string_id, ok := vars["user_id"]
	if !ok {
		http.Error(w, "missing user_id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(string_id)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	user, err := GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func EditUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
