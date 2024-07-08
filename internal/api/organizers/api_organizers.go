package organizers

import (
	mongodb "Account-Service/internal/database/mongodb/organizers"
	"Account-Service/internal/models"
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

func EditOrganizerById(w http.ResponseWriter, r *http.Request) {

	type RequestOrganizer struct {
		Photo       string `json:"photo"`
		Name        string `json:"name"`
		Description string `json:"description"`
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
