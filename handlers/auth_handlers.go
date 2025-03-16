package handlers

import (
	"encoding/json"
	"golang-practice/auth"
	"golang-practice/db"
	"golang-practice/models"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error": "Error hashing password"}`, http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	db.DB.Create(&user)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	var authReq AuthRequest
	json.NewDecoder(r.Body).Decode(&authReq)
	var user models.User
	db.DB.Where("email = ?", authReq.Email).First(&user)
	if user.Email == "" || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authReq.Password)) != nil {
		http.Error(w, `{"error": "Invalid credentials"}`, http.StatusUnauthorized)
		return
	}
	token, _ := auth.GenerateToken(user.Email, user.Role)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TokenResponse{Token: token})
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			http.Error(w, `{"error": "Missing token"}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func ProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Welcome to the protected route!"})
}
