package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"beads-haven/internal/httpjson"
	"beads-haven/internal/middleware"
	"beads-haven/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *sql.DB
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		httpjson.WriteError(w, http.StatusServiceUnavailable, "Database unavailable")
		return
	}

	var req models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Email == "" || req.Password == "" {
		httpjson.WriteError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	role := "user"
	var id int
	err = h.DB.QueryRow(`
		INSERT INTO users (email, password_hash, role)
		VALUES ($1, $2, $3)
		RETURNING id
	`, req.Email, string(hash), role).Scan(&id)

	if err != nil {
		httpjson.WriteError(w, http.StatusConflict, "Email already registered")
		return
	}

	token, err := middleware.GenerateToken(req.Email, role)
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.AuthResponse{
		Token: token,
		Role:  role,
		Email: req.Email,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		httpjson.WriteError(w, http.StatusServiceUnavailable, "Database unavailable")
		return
	}

	var req models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var user models.User
	err := h.DB.QueryRow(`
		SELECT id, email, password_hash, role, created_at
		FROM users
		WHERE email = $1
	`, req.Email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)

	if err != nil {
		httpjson.WriteError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		httpjson.WriteError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := middleware.GenerateToken(user.Email, user.Role)
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.AuthResponse{
		Token: token,
		Role:  user.Role,
		Email: user.Email,
	})
}
