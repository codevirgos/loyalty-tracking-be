package main

import (
	"Payback_BE/handlers"
	"Payback_BE/repo"
	"log"
	"net/http"
)

const port string = ":8080"

/*
func LookupNumber(w http.ResponseWriter, r *http.Request) {
	log.Println("looking up number..")

	if r.URL.Query().Has("nkey") {
		number := r.URL.Query().Get(("nkey"))
		fmt.Printf("number is %s \n", number)
	}
}

func RegisterNumber(w http.ResponseWriter, r *http.Request) {
	log.Println("registering number..")
}

func RedeemPoints(w http.ResponseWriter, r *http.Request) {
	log.Println("redeeming points")
}
*/

func main() {
	// repository pattern
	db, err := repo.InitPostgresDB()
	if err != nil {
		log.Printf("Failed to initialize db %v", err)
	}
	userHandler := handlers.NewUserHandler(db)
	http.HandleFunc("/lookup", userHandler.LookupNumber)
	http.HandleFunc("/register", userHandler.RegisterNumber)
	http.HandleFunc("/redeem", userHandler.RedeemPoints)

	log.Println("server started on port ", port)
	serveErr := http.ListenAndServe(port, nil)
	if serveErr != nil {
		log.Fatal(serveErr)
	}
}
