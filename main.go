package main

import (
	"fmt"
	"log"
	"net/http"

	"golang-practice/db"
	"golang-practice/handlers"

	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the JWT Auth API!")
}

func main() {
	db.InitDB()
	r := mux.NewRouter()

	// Add a home route to prevent 404 on root
	r.HandleFunc("/", HomeHandler).Methods("GET")

	// Auth Routes
	r.HandleFunc("/signup", handlers.SignUp).Methods("POST")
	r.HandleFunc("/signin", handlers.SignIn).Methods("POST")

	// Protected Route Example
	r.HandleFunc("/protected", handlers.AuthMiddleware(handlers.ProtectedEndpoint)).Methods("GET")

	// Start Server
	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", r)
}
