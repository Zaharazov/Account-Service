package services

import (
	"Account-Service/internal/domain/models"
	mongodb "Account-Service/internal/services/database/employers"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// GetAllEmployers godoc
// @Summary Get the list of all Employers
// @Description Получить список всех работодателей
// @Tags Employers
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Employer
// @Failure 400 {object} string "{"message"}"
// @Failure 404 {object} string "{"message"}"
// @Failure 500 {object} string "{"message"}"
// @Router /v1/users/employers [get]
func GetAllEmployers(w http.ResponseWriter, r *http.Request) {

	plug := uuid.New()
	employers, err := mongodb.GetEmployers(plug, 100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(employers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// GetEmployerById godoc
// @Summary Get Employer by id
// @Description Получить работодателя по его ID
// @Tags Employers
// @Accept  json
// @Produce  json
// @Param user_id path string true "User ID"
// @Success 200 {object} models.Employer
// @Failure 400 {object} string "{"message"}"
// @Failure 404 {object} string "{"message"}"
// @Failure 500 {object} string "{"message"}"
// @Router /v1/users/employers/{user_id} [get]
func GetEmployerById(w http.ResponseWriter, r *http.Request) {

	user_id, err := GetUserIdFromURL(w, r)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	employers, err := mongodb.GetEmployers(user_id, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	var employer models.Employer
	if len(employers) != 0 {
		employer = employers[0]
	} else {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(employer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// GetEmployersByParams godoc
// @Summary Find Employer by parameters
// @Description Найти определенных работодателей
// @Tags Employers
// @Accept  json
// @Produce  json
// @Param Employer body models.Employer true "Employer"
// @Success 200 {array} models.Employer
// @Failure 400 {object} string "{"message"}"
// @Failure 304 {object} string "{"message"}"
// @Router /v1/users/employers/search [get]
func GetEmployersByParams(w http.ResponseWriter, r *http.Request) {

	type RequestParams struct {
		Name        string `json:"name"`
		Focus       string `json:"focus"`
		Description string `json:"description"`
	}

	dec := json.NewDecoder(r.Body) // декодируем тело запроса
	dec.DisallowUnknownFields()    // проверяем, что получили то, что готовы принять (ругается на id, если оно есть в запросе)
	var rp RequestParams
	if err := dec.Decode(&rp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // получаем декодированные данные и проверяем, что все ок
		return
	}

	employers, err := mongodb.GetEmployersByParams(rp.Name, rp.Focus, rp.Description, 100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(employers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// EditEmployerById godoc
// @Summary Update Employer by id
// @Description Обновить существующего работодателя
// @Tags Employers
// @Accept  json
// @Produce  json
// @Param user_id path string true "User ID"
// @Param Employer body models.Employer true "Employer"
// @Success 200 {object} string "{"user_id"}"
// @Failure 400 {object} string "{"message"}"
// @Failure 304 {object} string "{"message"}"
// @Router /v1/users/employers/{user_id} [put]
func EditEmployerById(w http.ResponseWriter, r *http.Request) {

	type RequestEmployer struct {
		Photo       string `json:"photo,omitempty"`
		Name        string `json:"name,omitempty"`
		Focus       string `json:"focus,omitempty"`
		Description string `json:"description,omitempty"`
	}

	type ResponseEmployerId struct {
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
	var re RequestEmployer
	if err := dec.Decode(&re); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // получаем декодированные данные и проверяем, что все ок
		return
	}

	if len(re.Name) == 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := mongodb.EditEmployer(user_id, re.Photo, re.Name, re.Focus, re.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotModified)
		return
	}

	js, err := json.Marshal(ResponseEmployerId{Id: id}) // формируем json ответ с id выше
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
