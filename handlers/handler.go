package handlers

import (
	"Payback_BE/models"
	"Payback_BE/services"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type UserHandler struct {
	service *services.UserService
}

type Response struct {
	Code    int
	Data    *models.User
	Message string
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{
		service: services.NewUserService(db),
	}
}

// get
func (h *UserHandler) LookupNumber(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Has("nkey") {
		number := r.URL.Query().Get(("nkey"))
		num, err := strconv.Atoi(number)
		if err != nil {
			log.Printf("Failed to convert string to int %v", err)
		}
		user, lookupErr := h.service.LookupNumber(r.Context(), num, 10) // todo come back to this
		if lookupErr != nil {
			log.Println("Number not found") // thats okay

			res := &Response{
				Code:    204,
				Data:    &models.User{},
				Message: "number not found",
			}
			w.WriteHeader(204)
			json.NewEncoder(w).Encode(res)
		}
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

// post
func (h *UserHandler) RegisterNumber(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		log.Println("Failed to convert into obj")
	}

	user, registerErr := h.service.RegisterNumber(r.Context(), user.Number, 10) // todo come back to this
	if registerErr != nil {
		log.Println("Number not registered")

		res := &Response{
			Code:    500,
			Data:    &models.User{},
			Message: "number not registered",
		}
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(res)

	}
	res := &Response{
		Code:    200,
		Data:    user,
		Message: "number registered",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// put
func (h *UserHandler) RedeemPoints(w http.ResponseWriter, r *http.Request) {
	redeem := &models.Redeem{}

	err := json.NewDecoder(r.Body).Decode(redeem)
	if err != nil {
		log.Println("Failed to convert into obj")
	}

	redeemedUser, redeemedErr := h.service.RedeemPoints(r.Context(), redeem.Number, redeem.Amount)
	if redeemedErr != nil {
		log.Println("Points not redeemed")
		res := &Response{
			Code:    500,
			Data:    &models.User{},
			Message: "points not redeemed",
		}
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(res)
	}
	res := &Response{
		Code:    200,
		Data:    redeemedUser,
		Message: "points redeemed",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
