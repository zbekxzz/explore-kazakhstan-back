package handlers

import (
	"auth/internal/jwt"
	"auth/internal/repository"
	"auth/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

func RegisterHandler(repo *repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Хешируем пароль
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		// Создаем пользователя в базе данных
		err = repo.Create(user.Email, hashedPassword)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "User %s registered successfully", user.Email)
	}
}

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

		accessToken, refreshToken, err := jwt.GenerateTokens(user.Email)
		if err != nil {
			http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "Login successful", "accessToken": accessToken, "refresh_token": refreshToken})
	}
}

func ListUsersHandler(repo *repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

func UpdateUserHandler(repo *repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, err := jwt.ExtractUserEmail(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Декодируем новые данные пользователя
		var user struct {
			NewEmail    string `json:"new_email"`
			NewPassword string `json:"new_password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		// Хешируем новый пароль
		hashedPassword, _ := utils.HashPassword(user.NewPassword)

		// Обновляем данные в базе
		err = repo.UpdateByEmail(email, user.NewEmail, hashedPassword)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to update user: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "User with email %s updated", email)
	}
}

func GetUserHandler(repo *repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, err := jwt.ExtractUserEmail(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		user, err := repo.FindByEmail(email)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to retrieve user: %v", err), http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}
