package services

import (
	"Account-Service/internal/domain/models"
	mongodb "Account-Service/internal/services/database/students"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// GetAllStudents godoc
// @Summary Get the list of all Students
// @Description Получить список всех студентов
// @Tags Students
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Student
// @Failure 400 {object} string "{"message"}"
// @Failure 404 {object} string "{"message"}"
// @Failure 500 {object} string "{"message"}"
// @Router /v1/users/students [get]
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

// GetStudentById godoc
// @Summary Get Student by id
// @Description Получить студента по его ID
// @Tags Students
// @Accept  json
// @Produce  json
// @Param user_id path string true "User ID"
// @Success 200 {object} models.Student
// @Failure 400 {object} string "{"message"}"
// @Failure 404 {object} string "{"message"}"
// @Failure 500 {object} string "{"message"}"
// @Router /v1/users/students/{user_id} [get]
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

// GetStudentsByParams godoc
// @Summary Find Student by parameters
// @Description Найти определенных студентов
// @Tags Students
// @Accept  json
// @Produce  json
// @Param Student body models.Student true "Student"
// @Success 200 {array} models.Student
// @Failure 400 {object} string "{"message"}"
// @Failure 304 {object} string "{"message"}"
// @Router /v1/users/students/search [get]
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

// EditStudentById godoc
// @Summary Update Student by id
// @Description Обновить существующего студента
// @Tags Students
// @Accept  json
// @Produce  json
// @Param user_id path string true "User ID"
// @Param Student body models.Student true "Student"
// @Success 200 {object} string "{"user_id"}"
// @Failure 400 {object} string "{"message"}"
// @Failure 304 {object} string "{"message"}"
// @Router /v1/users/students/{user_id} [put]
func EditStudentById(w http.ResponseWriter, r *http.Request) {

	type RequestStudent struct {
		Photo       string `json:"photo,omitempty"`
		FullName    string `json:"full_name,omitempty"`
		Group       string `json:"group,omitempty"`
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
