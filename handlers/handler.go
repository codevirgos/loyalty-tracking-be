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
)

type UserHandler struct {
	UserTaskService *services.UserTaskService
}

type Response struct {
	Code    int
	Data    *models.User
	Message string
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{
		UserTaskService: services.NewUserTaskService(db),
	}
}

func (h *UserHandler) FindUser(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Has("id") {
		idString := r.URL.Query().Get(("id"))
		id, err := strconv.Atoi(idString)
		if err != nil {
			log.Printf("Failed to convert string to int %v", err)
		}
		user, serviceErr := h.UserTaskService.FindUser(r.Context(), id)

		if serviceErr != nil {
			log.Print("user not found")

			res := &Response{
				Code:    204,
				Data:    nil,
				Message: "user not found",
			}
			w.WriteHeader(204)
			json.NewEncoder(w).Encode(res)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			res := &Response{
				Code:    200,
				Data:    user,
				Message: "user found",
			}
			json.NewEncoder(w).Encode(res)
		}
	}
}

func (h *UserHandler) FindUserFromNum(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Has("nkey") {
		number := r.URL.Query().Get(("nkey"))
		user, serviceErr := h.UserTaskService.FindUserFromNum(r.Context(), number)

		// user not found
		if serviceErr != nil {
			log.Println("user not found") // thats okay

			res := &Response{
				Code:    204,
				Data:    nil,
				Message: "user not found",
			}
			w.WriteHeader(204)
			json.NewEncoder(w).Encode(res)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			res := &Response{
				Code:    200,
				Data:    user,
				Message: "user found",
			}
			json.NewEncoder(w).Encode(res)
		}
	}
}

func (h *UserHandler) SaveUser(w http.ResponseWriter, r *http.Request) {
	var u models.User

	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to convert string to int %v", err)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(b, &u); err != nil {
		http.Error(w, "Conversion error", http.StatusBadRequest)
		return
	}

	savedUser, serviceErr := h.UserTaskService.SaveUser(r.Context(), u.Number, u.Name, u.ID)
	if serviceErr != nil {
		log.Println("user not saved")
		res := &Response{
			Code:    204,
			Data:    nil,
			Message: "user not saved",
		}
		w.WriteHeader(204)
		json.NewEncoder(w).Encode(res)
	} else {
		log.Println("user saved")
		res := &Response{
			Code:    200,
			Data:    savedUser,
			Message: "user saved",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}
