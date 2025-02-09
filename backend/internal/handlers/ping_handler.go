package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"backend/internal/domain"
	"backend/internal/usecase"

	"github.com/gorilla/mux"
)

type PingHandler struct {
	pingUsecase usecase.PingUsecase
}

func NewPingHandler(pingUsecase usecase.PingUsecase) *PingHandler {
	return &PingHandler{pingUsecase: pingUsecase}
}

func (h *PingHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/pings", h.GetPings).Methods("GET")
	router.HandleFunc("/pings/add", h.AddPing).Methods("POST")
}

func (h *PingHandler) GetPings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pings, err := h.pingUsecase.GetAllPings()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(pings)
}

func (h *PingHandler) AddPing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Разрешен только метод POST", http.StatusMethodNotAllowed)
		return
	}

	ipAddress := r.URL.Query().Get("ip_address")
	pingTimeStr := r.URL.Query().Get("ping_time")
	if ipAddress == "" || pingTimeStr == "" {
		http.Error(w, "Требуются значения ip_address и ping_time", http.StatusBadRequest)
		return
	}

	pingTime, err := strconv.Atoi(pingTimeStr)
	if err != nil {
		http.Error(w, "Недопустимое значение ping_time", http.StatusBadRequest)
		return
	}

	pingData := domain.PingData{
		IPAddress:          ipAddress,
		PingTime:           pingTime,
		LastSuccessfulPing: time.Now(),
	}

	err = h.pingUsecase.AddPing(pingData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Ping добавлен успешно")
}
