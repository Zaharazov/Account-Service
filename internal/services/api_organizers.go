package services

import (
	"Account-Service/internal/domain/models"
	mongodb "Account-Service/internal/services/database/organizers"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetEventIdFromURL(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	vars := mux.Vars(r)
	string_id, ok := vars["event_id"]
	if !ok {
		http.Error(w, "missing event_id", http.StatusBadRequest)
	}
	event_id, err := uuid.Parse(string_id)
	if err != nil {
		return event_id, err
	}
	return event_id, nil
}

func AddEventToOrganizerById(w http.ResponseWriter, r *http.Request) { // PUT

	_ = ContentTypeCheck(w, r)

	user_id, err := GetUserIdFromURL(w, r)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	event_id, err := GetEventIdFromURL(w, r)
	if err != nil {
		http.Error(w, "invalid event_id", http.StatusBadRequest)
		return
	}

	err = mongodb.AddEvent(user_id, event_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotModified)
		return
	}

	organizers, err := mongodb.GetOrganizers(user_id, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var organizer models.Organizer
	if len(organizers) != 0 {
		organizer = organizers[0]
	} else {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(organizer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// GetAllOrganizers godoc
// @Summary Get the list of all Organizers
// @Description Получить список всех организаторов
// @Tags Organizers
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Organizer
// @Failure 400 {object} string "{"message"}"
// @Failure 404 {object} string "{"message"}"
// @Failure 500 {object} string "{"message"}"
// @Router /v1/users/organizers [get]
func GetAllOrganizers(w http.ResponseWriter, r *http.Request) {

	plug := uuid.New()
	organizers, err := mongodb.GetOrganizers(plug, 100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(organizers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// GetOrganizerById godoc
// @Summary Get Organizer by id
// @Description Получить организатора по его ID
// @Tags Organizers
// @Accept  json
// @Produce  json
// @Param user_id path string true "User ID"
// @Success 200 {object} models.Organizer
// @Failure 400 {object} string "{"message"}"
// @Failure 404 {object} string "{"message"}"
// @Failure 500 {object} string "{"message"}"
// @Router /v1/users/organizers/{user_id} [get]
func GetOrganizerById(w http.ResponseWriter, r *http.Request) {

	user_id, err := GetUserIdFromURL(w, r)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	organizers, err := mongodb.GetOrganizers(user_id, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	var organizer models.Organizer
	if len(organizers) != 0 {
		organizer = organizers[0]
	} else {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(organizer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// GetOrganizersByParams godoc
// @Summary Find Organizer by parameters
// @Description Найти определенных организаторов
// @Tags Organizers
// @Accept  json
// @Produce  json
// @Param Organizer body models.Organizer true "Organizer"
// @Success 200 {array} models.Organizer
// @Failure 400 {object} string "{"message"}"
// @Failure 304 {object} string "{"message"}"
// @Router /v1/users/organizers/search [get]
func GetOrganizersByParams(w http.ResponseWriter, r *http.Request) {

	type RequestParams struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	dec := json.NewDecoder(r.Body) // декодируем тело запроса
	dec.DisallowUnknownFields()    // проверяем, что получили то, что готовы принять (ругается на id, если оно есть в запросе)
	var rp RequestParams
	if err := dec.Decode(&rp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // получаем декодированные данные и проверяем, что все ок
		return
	}

	organizers, err := mongodb.GetOrganizersByParams(rp.Name, rp.Description, 100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(organizers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// EditOrganizerById godoc
// @Summary Update Organizer by id
// @Description Обновить существующего организатора
// @Tags Organizers
// @Accept  json
// @Produce  json
// @Param user_id path string true "User ID"
// @Param Organizer body models.Organizer true "Organizer"
// @Success 200 {object} string "{"user_id"}"
// @Failure 400 {object} string "{"message"}"
// @Failure 304 {object} string "{"message"}"
// @Router /v1/users/organizers/{user_id} [put]
func EditOrganizerById(w http.ResponseWriter, r *http.Request) {

	type RequestOrganizer struct {
		Photo       string `json:"photo,omitempty"`
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
	}

	type ResponseOrganizerId struct {
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
	var ro RequestOrganizer
	if err := dec.Decode(&ro); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // получаем декодированные данные и проверяем, что все ок
		return
	}

	if len(ro.Name) == 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := mongodb.EditOrganizer(user_id, ro.Photo, ro.Name, ro.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotModified)
		return
	}

	js, err := json.Marshal(ResponseOrganizerId{Id: id}) // формируем json ответ с id выше
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
