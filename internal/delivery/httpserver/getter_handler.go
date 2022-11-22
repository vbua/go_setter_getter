package httpserver

import (
	"github.com/gorilla/mux"
	"github.com/vbua/go_setter_getter/internal/entity"
	"net/http"
)

type UserGradeGetterService interface {
	Get(userId string) (*entity.UserGrade, error)
}

type GetterHandler struct {
	userGradeService UserGradeGetterService
}

func NewGetterHandler(userGradeService UserGradeGetterService) GetterHandler {
	return GetterHandler{userGradeService}
}

func (h *GetterHandler) CreateRoutes() *mux.Router {
	getter := mux.NewRouter()
	getter.HandleFunc("/get", h.get).Methods("GET")
	getter.Use(logging)
	return getter
}

func (h *GetterHandler) get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := r.URL.Query().Get("user_id")
	if userId == "" {
		sendResponse(true, "user id is required", http.StatusBadRequest, w)
		return
	}
	userGrade, err := h.userGradeService.Get(userId)
	if err != nil {
		sendResponse(true, err.Error(), http.StatusInternalServerError, w)
		return
	}
	sendResponse(false, userGrade, http.StatusOK, w)
}
