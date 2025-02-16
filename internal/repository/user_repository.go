package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

type User struct {
	ID       int
	Email    string
	Password string
	Role     string
}

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) FindByEmail(email string) (*User, error) {
	var user User
	err := repo.DB.QueryRow("SELECT id, email, password, role FROM users WHERE email=$1", email).
		Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Пользователь не найден
		}
		return nil, err
	}
	return &user, nil
}

// GetAllUsers возвращает список всех пользователей
func (repo *UserRepository) GetAllUsers() ([]User, error) {
	rows, err := repo.DB.Query("SELECT id, email, role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Email, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *UserRepository) Create(email, password, role string) error {
	_, err := repo.DB.Exec("INSERT INTO users (email, password, role) VALUES ($1, $2, $3)", email, password, role)
	return err
}

func (repo *UserRepository) UpdateByID(id string, newEmail, newRole string) error {
	userID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	_, err = repo.DB.Exec("UPDATE users SET email=$1, role=$2 WHERE id=$3", newEmail, newRole, userID)
	return err
}

func (repo *UserRepository) DeleteByID(id string) error {
	userID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	_, err = repo.DB.Exec("DELETE FROM users WHERE id=$1", userID)
	return err
}
