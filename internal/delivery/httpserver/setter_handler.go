package httpserver

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/vbua/go_setter_getter/internal/entity"
	"net/http"
)

type UserGradeSetterService interface {
	Set(entity.UserGrade)
}

type SetterHandler struct {
	userGradeService UserGradeSetterService
}

func NewSetterHandler(userGradeService UserGradeSetterService) SetterHandler {
	return SetterHandler{userGradeService}
}

func (h *SetterHandler) CreateRoutes() *mux.Router {
	setter := mux.NewRouter()
	setter.HandleFunc("/set", h.set).Methods("POST")
	setter.Use(logging)
	setter.Use(basicAuth)
	return setter
}

func (h *SetterHandler) set(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userGrade := entity.UserGrade{}
	if err := json.NewDecoder(r.Body).Decode(&userGrade); err != nil {
		sendResponse(true, err.Error(), http.StatusInternalServerError, w)
		return
	}
	validate := validator.New()
	err := validate.Struct(userGrade)
	if err != nil {
		sendResponse(true, err.Error(), http.StatusBadRequest, w)
		return
	}
	h.userGradeService.Set(userGrade)
	sendResponse(false, "User grade created successfully", http.StatusOK, w)
}
