package services

import (
	"Account-Service/internal/domain/models"
	mongodb_e "Account-Service/internal/services/database/employers"
	mongodb_o "Account-Service/internal/services/database/organizers"
	mongodb_s "Account-Service/internal/services/database/students"
	mongodb_u "Account-Service/internal/services/database/users"
	"encoding/json"
	"mime"
	"net/http"

	"github.com/google/uuid"
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

func GetUserIdFromURL(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	vars := mux.Vars(r)
	string_id, ok := vars["user_id"]
	if !ok {
		http.Error(w, "missing user_id", http.StatusBadRequest)
	}
	user_id, err := uuid.Parse(string_id)
	if err != nil {
		return user_id, err
	}
	return user_id, nil
}

// CreateUser godoc
// @Summary Create new User
// @Description Создать нового пользователя
// @Tags Users
// @Accept  json
// @Produce  json
// @Param User body models.User true "User need to be created"
// @Success 201 {object} string "{"user_id"}"
// @Failure 400 {object} string "{"message"}"
// @Failure 500 {object} string "{"message"}"
// @Router /v1/users/ [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {

	type RequestUser struct {
		Login    string   `json:"login,omitempty"`
		Password string   `json:"password,omitempty"`
		Roles    []string `json:"roles,omitempty"`
	}

	type ResponseUserId struct {
		Id uuid.UUID `json:"id"` // структура объекта, который отдаем
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

	user_id := uuid.New()
	id, err := mongodb_u.SaveUser(user_id, ru.Login, ru.Password, ru.Roles) // создаем карточку по данным из запроса и получаем его id

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, role := range ru.Roles {
		if role == "student" {
			_, err = mongodb_s.SaveStudent(user_id, "", "", "", "", "", "", "")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		if role == "organizer" {
			_, err = mongodb_o.SaveOrganizer(user_id, "", "", "")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		if role == "employer" {
			_, err = mongodb_e.SaveEmployer(user_id, "", "", "", "")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

	}

	js, err := json.Marshal(ResponseUserId{Id: id}) // формируем json ответ с id выше
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(js)
}

// DeleteUserById godoc
// @Summary Delete User by id
// @Description Удалить пользователя по ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user_id path string true "User ID"
// @Success 200
// @Failure 400 {object} string "{"message"}"
// @Failure 404 {object} string "{"message"}"
// @Failure 500 {object} string "{"message"}"
// @Router /v1/users/{user_id} [delete]
func DeleteUserById(w http.ResponseWriter, r *http.Request) {

	user_id, err := GetUserIdFromURL(w, r)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	err = mongodb_u.DeleteUser(user_id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	mongodb_s.DeleteStudent(user_id)
	mongodb_o.DeleteOrganizer(user_id)
	mongodb_e.DeleteEmployer(user_id)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// GetAllUsers godoc
// @Summary Get the list of all Users
// @Description Получить список всех пользователей
// @Tags Users
// @Accept  json
// @Produce  json
// @Success 200 {array} models.User
// @Failure 400 {object} string "{"message"}"
// @Failure 404 {object} string "{"message"}"
// @Failure 500 {object} string "{"message"}"
// @Router /v1/users [get]
func GetAllUsers(w http.ResponseWriter, r *http.Request) {

	plug := uuid.New()
	users, err := mongodb_u.GetUsers(plug, 100)
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

// GetUserById godoc
// @Summary Get User by id
// @Description Получить пользователя по его ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user_id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} string "{"message"}"
// @Failure 404 {object} string "{"message"}"
// @Failure 500 {object} string "{"message"}"
// @Router /v1/users/{user_id} [get]
func GetUserById(w http.ResponseWriter, r *http.Request) {

	user_id, err := GetUserIdFromURL(w, r)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	users, err := mongodb_u.GetUsers(user_id, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	var user models.User
	if len(users) != 0 {
		user = users[0]
	} else {
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

// EditUserById godoc
// @Summary Update User by id
// @Description Обновить существующего пользователя
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user_id path string true "User ID"
// @Param User body models.User true "User"
// @Success 200 {object} string "{"user_id"}"
// @Failure 400 {object} string "{"message"}"
// @Failure 304 {object} string "{"message"}"
// @Router /v1/users/{user_id} [put]
func EditUserById(w http.ResponseWriter, r *http.Request) {

	type RequestUser struct {
		Login    string   `json:"login,omitempty"`
		Password string   `json:"password,omitempty"`
		Roles    []string `json:"roles,omitempty"`
	}

	type ResponseUserId struct {
		Id uuid.UUID `json:"id"` // структура объекта, который отдаем
	}

	_ = ContentTypeCheck(w, r)

	user_id, err := GetUserIdFromURL(w, r)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
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

	id, err := mongodb_u.EditUser(user_id, ru.Login, ru.Password, ru.Roles)
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
