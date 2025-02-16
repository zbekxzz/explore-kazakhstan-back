package handlers

import (
	"auth/internal/jwt"
	"auth/internal/repository"
	"auth/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

func LoginHandler(repo *repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		user, err := repo.FindByEmail(creds.Email)
		if err != nil {
			http.Error(w, "Error retrieving user", http.StatusInternalServerError)
			return
		}
		if user == nil || !utils.CheckPasswordHash(creds.Password, user.Password) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		accessToken, refreshToken, err := jwt.GenerateTokens(user.Email, user.Role)
		if err != nil {
			http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "Login successful", "accessToken": accessToken, "refresh_token": refreshToken})
	}
}

func ListUsersHandler(repo *repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := jwt.ExtractAdminClaims(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		users, err := repo.GetAllUsers()
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to retrieve users: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, "failed to encode users", http.StatusInternalServerError)
		}
	}
}

func CreateUserHandler(repo *repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := jwt.ExtractAdminClaims(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		var user struct {
			Email string `json:"email"`
			Role  string `json:"role"`
		}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		tempPassword := "Temp1234"
		hashedPassword, _ := utils.HashPassword(tempPassword)

		err = repo.Create(user.Email, hashedPassword, user.Role)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to create user: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "User %s created with temporary password: %s", user.Email, tempPassword)
	}
}

func UpdateUserHandler(repo *repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := jwt.ExtractAdminClaims(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "missing user ID in query parameters", http.StatusBadRequest)
			return
		}

		var user struct {
			NewRole  string `json:"new_role"`  // Новая роль
			NewEmail string `json:"new_email"` // Новый email
		}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		err = repo.UpdateByID(id, user.NewEmail, user.NewRole)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to update user: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "User with ID %s updated", id)
	}
}

func DeleteUserHandler(repo *repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := jwt.ExtractAdminClaims(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "missing user ID in query parameters", http.StatusBadRequest)
			return
		}

		err = repo.DeleteByID(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to delete user: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "User with ID %s deleted", id)
	}
}
