package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
	"example.com/go-rest/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *sql.DB
}

type AuthRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

var jwtKey = []byte("kunci_rahasia_super_aman_anda")

func (handler *AuthHandler) Register(writer http.ResponseWriter, request *http.Request){
	var authReq AuthRequest

	if err := json.NewDecoder(request.Body).Decode(&authReq); err != nil {
		http.Error(writer, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	if authReq.Username == "" || authReq.Email == "" || authReq.Password == "" {
		http.Error(writer, "Username, email, and password are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(authReq.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(writer, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	_, err = handler.DB.Exec("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)", authReq.Username, authReq.Email, string(hashedPassword))

	if err != nil {
		http.Error(writer, "Failed to create user. Email/Username may already exist.", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(map[string]string{"message": "User created successfully"})
}

func (handler *AuthHandler) Login(writer http.ResponseWriter, request *http.Request) {
	var authReq AuthRequest

	if err := json.NewDecoder(request.Body).Decode(&authReq); err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user models.User

	err := handler.DB.QueryRow("SELECT id, username, password_hash FROM users WHERE email = ?", authReq.Email).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(writer, "Invalid credentials", http.StatusUnauthorized)
		} else {
			http.Error(writer, "Internal Server Errro", http.StatusInternalServerError)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(authReq.Password)); err != nil {
		http.Error(writer, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Subject: user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		http.Error(writer, "Failed to create token", http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(AuthResponse{Token: tokenString})
}
