package users

import (
	"Account-Service/internal/api/users/mongo"
	"encoding/json"
	"mime"
	"net/http"

	"strconv"

	"github.com/gorilla/mux"
)

func ContentTypeCheck(w http.ResponseWriter, r *http.Request) error {
	contentType := r.Header.Get("Content-Type")           // получаем тип контента в запросе
	mediatype, _, err := mime.ParseMediaType(contentType) // парсинг полученных данных из запроса
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // ошибка, если не получили запрос
		return err
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType) // ошибка, если пришел не json
		return err
	}
	return nil
}

func GetUserIdFromURL(w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	string_id, ok := vars["user_id"]
	if !ok {
		http.Error(w, "missing user_id", http.StatusBadRequest)
	}
	user_id, err := strconv.Atoi(string_id)
	if err != nil {
		return -1, err
	}
	return user_id, nil
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

	err := ContentTypeCheck(w, r)

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

	id := mongo.SaveUser(ru.Login, ru.Password, ru.Roles) // создаем карточку по данным из запроса и получаем его id
	js, err := json.Marshal(ResponseUserId{Id: id})       // формируем json ответ с id выше
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func DeleteUserById(w http.ResponseWriter, r *http.Request) {

	user_id, err := GetUserIdFromURL(w, r)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
	}

	err = mongo.DeleteUser(user_id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {

	users, err := mongo.GetUsers(-1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {

	user_id, err := GetUserIdFromURL(w, r)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
	}

	users, err := mongo.GetUsers(user_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	user := users[0]
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

	type RequestUser struct {
		Login    string   `json:"login,omitempty"`
		Password string   `json:"password,omitempty"`
		Roles    []string `json:"roles,omitempty"`
	}

	type ResponseUserId struct {
		Id int32 `json:"id"` // структура объекта, который отдаем
	}

	_ = ContentTypeCheck(w, r)

	user_id, err := GetUserIdFromURL(w, r)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
	}

	dec := json.NewDecoder(r.Body) // декодируем тело запроса
	dec.DisallowUnknownFields()    // проверка, что получили то, что готовы принять (ругается на id, если оно есть в запросе)
	var ru RequestUser
	if err := dec.Decode(&ru); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // получаем декодированные данные и проверяем, что все ок
		return
	}
	id, err := mongo.EditUser(user_id, ru.Login, ru.Password, ru.Roles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotModified)
		return
	}

	js, err := json.Marshal(ResponseUserId{Id: id}) // формируем json ответ с id выше
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
