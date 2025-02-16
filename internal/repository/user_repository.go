package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

type User struct {
	ID       int     `json:"id"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Events   []Event `json:"events"`
}

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) Create(email, password string) error {
	_, err := repo.DB.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", email, password)
	return err
}

func (repo *UserRepository) FindByEmail(email string) (*User, error) {
	var user User
	err := repo.DB.QueryRow("SELECT id, email, password FROM users WHERE email=$1", email).
		Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	rows, err := repo.DB.Query(`
		SELECT id, user_id, title, date, time, venue, description, note, price, image_url, attendees, is_active, created_at
		FROM events WHERE user_id=$1`, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event Event
		var attendeesArray []string
		if err := rows.Scan(&event.ID, &event.UserID, &event.Title, &event.Date, &event.Time, &event.Venue,
			&event.Description, &event.Note, &event.Price, &event.ImageURL, &attendeesArray,
			&event.IsActive, &event.CreatedAt); err != nil {
			return nil, err
		}
		event.Attendees = attendeesArray
		user.Events = append(user.Events, event)
	}

	return &user, nil
}

func (repo *UserRepository) GetAllUsers() ([]User, error) {
	rows, err := repo.DB.Query("SELECT id, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Email); err != nil {
			return nil, err
		}

		eventRows, err := repo.DB.Query(`
			SELECT id, user_id, title, date, time, venue, description, note, price, image_url, attendees, is_active, created_at 
			FROM events WHERE user_id=$1`, user.ID)
		if err != nil {
			return nil, err
		}
		defer eventRows.Close()

		for eventRows.Next() {
			var event Event
			var attendeesArray []string
			if err := eventRows.Scan(&event.ID, &event.UserID, &event.Title, &event.Date, &event.Time, &event.Venue,
				&event.Description, &event.Note, &event.Price, &event.ImageURL, &attendeesArray,
				&event.IsActive, &event.CreatedAt); err != nil {
				return nil, err
			}
			event.Attendees = attendeesArray
			user.Events = append(user.Events, event)
		}

		users = append(users, user)
	}

	return users, nil
}

func (repo *UserRepository) UpdateByID(id string, newEmail, newPassword string) error {
	userID, err := strconv.Atoi(id)

	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	_, err = repo.DB.Exec("UPDATE users SET email=$1, password=$2 WHERE id=$3", newEmail, newPassword, userID)
	return err
}

func (repo *UserRepository) DeleteByID(id string) error {
	userID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	_, err = repo.DB.Exec("DELETE FROM events WHERE user_id=$1", userID)
	if err != nil {
		return err
	}

	_, err = repo.DB.Exec("DELETE FROM users WHERE id=$1", userID)
	return err
}
