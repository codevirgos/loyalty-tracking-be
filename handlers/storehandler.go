package handlers

import (
	"Payback_BE/models"
	"Payback_BE/services"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/skip2/go-qrcode"
)

type StoreHandler struct {
	StoreTaskService *services.StoreTaskService
}

type OverviewResponse struct {
	Code    int
	Data    *models.Overview
	Message string
}

func NewStoreHandler(db *sql.DB) *StoreHandler {
	return &StoreHandler{
		StoreTaskService: services.NewStoreTaskService(db),
	}
}

func (h *StoreHandler) GetOverview(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Has("id") {
		idString := r.URL.Query().Get(("id"))
		id, err := strconv.Atoi(idString)
		if err != nil {
			log.Printf("Failed to convert string to int %v", err)
		}
		overview, serviceErr := h.StoreTaskService.GetOverview(r.Context(), id)

		if serviceErr != nil {
			log.Print("no overview")

			res := &OverviewResponse{
				Code:    204,
				Data:    nil,
				Message: "no overview",
			}
			w.WriteHeader(204)
			json.NewEncoder(w).Encode(res)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			res := &OverviewResponse{
				Code:    200,
				Data:    overview,
				Message: "overview found",
			}
			json.NewEncoder(w).Encode(res)
		}
	}
}

func (h *StoreHandler) StoreScan(w http.ResponseWriter, r *http.Request) {
	var s models.Scanned
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body%v", err)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(b, &s); err != nil {
		http.Error(w, "Conversion error", http.StatusBadRequest)
		return
	}
	overview, serviceErr := h.StoreTaskService.StoreScan(r.Context(), s)

	if serviceErr != nil {
		log.Print("scan error")
		res := &OverviewResponse{
			Code:    500,
			Data:    nil,
			Message: "scan error",
		}
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(res)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res := &OverviewResponse{
			Code:    200,
			Data:    overview,
			Message: "scan successful",
		}
		json.NewEncoder(w).Encode(res)
	}
}

func (h *StoreHandler) RedeemPoints(w http.ResponseWriter, r *http.Request) {
	var rd models.Redeem
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body%v%", err)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(b, &rd); err != nil {
		http.Error(w, "Conversion error", http.StatusBadRequest)
		return
	}

	overview, serviceErr := h.StoreTaskService.RedeemPoints(r.Context(), rd)

	if serviceErr != nil {
		log.Print("redeem error")
		res := &OverviewResponse{
			Code:    500,
			Data:    nil,
			Message: "redeem error",
		}
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(res)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		log.Print("redeem successful")
		res := &OverviewResponse{
			Code:    200,
			Data:    overview,
			Message: "redeem successful",
		}
		json.NewEncoder(w).Encode(res)
	}
}

func (h *StoreHandler) GenerateQRCode(w http.ResponseWriter, r *http.Request) {
	err := qrcode.WriteFile("{'storeID': '1'}", qrcode.Medium, 256, "qrcode.png")
	if err != nil {
		log.Print(err)
	}
}
