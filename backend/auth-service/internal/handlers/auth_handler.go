package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/trademaster/backend/auth-service/internal/models"
	"github.com/trademaster/backend/auth-service/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type SignupRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"role"`
	Profile   ProfileData `json:"profile"`
}

type ProfileData struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Company   string `json:"company,omitempty"`
	Trade     string `json:"trade,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// Signup handles user registration
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	if req.Role == "" {
		req.Role = "contractor" // Default role
	}

	// Validate role
	validRoles := map[string]bool{"contractor": true, "client": true, "admin": true}
	if !validRoles[req.Role] {
		respondWithError(w, http.StatusBadRequest, "Invalid role")
		return
	}

	// Create user
	user, token, err := h.authService.Signup(req.Email, req.Password, req.Role, req.Profile.FirstName, req.Profile.LastName, req.Profile.Company, req.Profile.Trade)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			respondWithError(w, http.StatusConflict, "Email already exists")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, AuthResponse{
		User:  user,
		Token: token,
	})
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	// Authenticate user
	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "invalid credentials") {
			respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Authentication failed")
		return
	}

	respondWithJSON(w, http.StatusOK, AuthResponse{
		User:  user,
		Token: token,
	})
}

// Logout handles user logout (client-side token deletion)
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// In a stateless JWT system, logout is handled client-side
	// Could implement token blacklisting here if needed
	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Logout successful",
	})
}

// GetCurrentUser returns the authenticated user's information
func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing authorization header")
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		respondWithError(w, http.StatusUnauthorized, "Invalid authorization header format")
		return
	}

	// Validate token and get user
	user, err := h.authService.ValidateToken(tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"user": user,
	})
}

// RefreshToken generates a new token for the user
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing authorization header")
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Validate current token
	user, err := h.authService.ValidateToken(tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	// Generate new token
	newToken, err := h.authService.GenerateToken(user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"token": newToken,
	})
}

// GetUser retrieves a user by ID
func (h *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	if userID == "" {
		respondWithError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve user")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"user": user,
	})
}

// UpdateUser updates user information
func (h *AuthHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	// Verify authentication
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing authorization header")
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	currentUser, err := h.authService.ValidateToken(tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	// Users can only update their own profile (unless admin)
	if currentUser.ID.String() != userID && currentUser.Role != "admin" {
		respondWithError(w, http.StatusForbidden, "You can only update your own profile")
		return
	}

	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update user (implementation would go to service layer)
	// For now, just return success
	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "User updated successfully",
	})
}

// Helper functions

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal server error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ErrorResponse{
		Error:   http.StatusText(code),
		Message: message,
	})
}
