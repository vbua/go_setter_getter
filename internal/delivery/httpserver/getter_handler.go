package httpserver

import (
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/vbua/go_setter_getter/internal/entity"
	"net/http"
	"strconv"
	"time"
)

type UserGradeGetterService interface {
	Get(userId string) (*entity.UserGrade, error)
	Backup() map[string]entity.UserGrade
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
	getter.HandleFunc("/backup", h.backup).Methods("GET")
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

func (h *GetterHandler) backup(w http.ResponseWriter, r *http.Request) {
	const layout = "01-02-2006 15:04:05"
	t := time.Now()
	w.Header().Set("Content-Disposition", `attachment; filename="`+t.Format(layout)+`.csv.gz"`)
	var buf bytes.Buffer
	userGrades := h.userGradeService.Backup()
	zipWriter := gzip.NewWriter(&buf)
	csvwriter := csv.NewWriter(zipWriter)

	for _, grade := range userGrades {
		r := make([]string, 0, 5)
		r = append(
			r,
			grade.UserId,
			strconv.Itoa(grade.PostpaidLimit),
			strconv.Itoa(grade.Spp),
			strconv.Itoa(grade.ShippingFee),
			strconv.Itoa(grade.ReturnFee),
		)
		csvwriter.Write(r)
	}
	csvwriter.Flush()
	zipWriter.Flush()
	zipWriter.Close()

	fmt.Println("Compressed data:", buf.Bytes())

	w.Write(buf.Bytes())
}
