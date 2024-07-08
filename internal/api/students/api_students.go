package students

import (
	mongodb "Account-Service/internal/database/mongodb/students"
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

func GetAllStudents(w http.ResponseWriter, r *http.Request) {

	plug := uuid.New()
	students, err := mongodb.GetStudents(plug, 100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(students)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func GetStudentById(w http.ResponseWriter, r *http.Request) {

	user_id, err := GetUserIdFromURL(w, r)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	students, err := mongodb.GetStudents(user_id, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	var student models.Student
	if len(students) != 0 {
		student = students[0]
	} else {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(student)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func GetStudentsByParams(w http.ResponseWriter, r *http.Request) {

	type RequestParams struct {
		FullName    string `json:"full_name"`
		Group       string `json:"group"`
		RecordBook  string `json:"record_book,omitempty"`
		Description string `json:"description,omitempty"`
		Mail        string `json:"mail,omitempty"`
		GitHub      string `json:"github,omitempty"`
	}

	dec := json.NewDecoder(r.Body) // декодируем тело запроса
	dec.DisallowUnknownFields()    // проверка, что получили то, что готовы принять (ругается на id, если оно есть в запросе)
	var rp RequestParams
	if err := dec.Decode(&rp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // получаем декодированные данные и проверяем, что все ок
		return
	}

	students, err := mongodb.GetStudentsByParams(rp.FullName, rp.Group, rp.RecordBook, rp.Description, rp.Mail, rp.GitHub, 100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(students)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func EditStudentById(w http.ResponseWriter, r *http.Request) {

	type RequestStudent struct {
		Photo       string `json:"photo,omitempty"`
		FullName    string `json:"full_name"`
		Group       string `json:"group"`
		RecordBook  string `json:"record_book,omitempty"`
		Description string `json:"description,omitempty"`
		Mail        string `json:"mail,omitempty"`
		GitHub      string `json:"github,omitempty"`
	}

	type ResponseStudentId struct {
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
	var rs RequestStudent
	if err := dec.Decode(&rs); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // получаем декодированные данные и проверяем, что все ок
		return
	}

	if len(rs.FullName) == 0 || len(rs.Group) == 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := mongodb.EditStudent(user_id, rs.Photo, rs.FullName, rs.Group, rs.RecordBook, rs.Description, rs.Mail, rs.GitHub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotModified)
		return
	}

	js, err := json.Marshal(ResponseStudentId{Id: id}) // формируем json ответ с id выше
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
