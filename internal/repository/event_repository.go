package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"strconv"
	"time"
)

type Event struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Date        string    `json:"date"`
	Time        string    `json:"time"`
	Venue       string    `json:"venue"`
	Description string    `json:"description"`
	Note        string    `json:"note"`
	Price       string    `json:"price"`
	ImageURL    string    `json:"image_url"`
	Attendees   []string  `json:"attendees"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type EventRepository struct {
	DB *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{DB: db}
}

func (repo *EventRepository) FindByID(id int) (*Event, error) {
	var event Event
	var attendeesArray []string

	err := repo.DB.QueryRow(`
		SELECT id, user_id, title, date, time, venue, description, note, price, image_url, attendees, is_active, created_at
		FROM events WHERE id=$1`, id).
		Scan(&event.ID, &event.UserID, &event.Title, &event.Date, &event.Time, &event.Venue,
			&event.Description, &event.Note, &event.Price, &event.ImageURL, pq.Array(&attendeesArray), // ✅ Читаем массив
			&event.IsActive, &event.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	event.Attendees = attendeesArray
	return &event, nil
}

func (repo *EventRepository) GetAllEvents() ([]Event, error) {
	rows, err := repo.DB.Query("SELECT id, user_id, title, date, time, venue, description, note, price, image_url, attendees, is_active, created_at FROM events")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		var attendeesArray []string
		if err := rows.Scan(&event.ID, &event.UserID, &event.Title, &event.Date, &event.Time,
			&event.Venue, &event.Description, &event.Note, &event.Price, &event.ImageURL,
			pq.Array(&attendeesArray), // ✅ Читаем массив `attendees`
			&event.IsActive, &event.CreatedAt); err != nil {
			return nil, err
		}
		event.Attendees = attendeesArray
		events = append(events, event)
	}

	return events, nil
}

func (repo *EventRepository) Create(event Event) error {
	_, err := repo.DB.Exec(`
		INSERT INTO events (user_id, title, date, time, venue, description, note, price, image_url, attendees, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		event.UserID, event.Title, event.Date, event.Time, event.Venue,
		event.Description, event.Note, event.Price, event.ImageURL,
		pq.Array(event.Attendees),
		event.IsActive)
	return err
}

func (repo *EventRepository) UpdateByID(id string, event Event) error {
	eventID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid event ID: %w", err)
	}

	_, err = repo.DB.Exec(`
		UPDATE events 
		SET title=$1, date=$2, time=$3, venue=$4, description=$5, note=$6, 
		    price=$7, image_url=$8, attendees=$9, is_active=$10
		WHERE id=$11`,
		event.Title, event.Date, event.Time, event.Venue, event.Description,
		event.Note, event.Price, event.ImageURL, pq.Array(event.Attendees), // ✅ Используем `pq.Array`
		event.IsActive, eventID)
	return err
}

func (repo *EventRepository) DeleteByID(id string) error {
	eventID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid event ID: %w", err)
	}

	_, err = repo.DB.Exec("DELETE FROM events WHERE id=$1", eventID)
	return err
}

func (repo *EventRepository) RegisterUserForEvent(eventID, userEmail string) (*Event, error) {
	var attendees []string
	var event Event

	err := repo.DB.QueryRow(`
		SELECT id, user_id, title, date, time, venue, description 
		FROM events WHERE id=$1`, eventID).
		Scan(&event.ID, &event.UserID, &event.Title, &event.Date, &event.Time, &event.Venue, &event.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("event not found")
		}
		return nil, err
	}

	err = repo.DB.QueryRow("SELECT attendees FROM events WHERE id=$1", eventID).Scan(pq.Array(&attendees))
	if err != nil {
		return nil, err
	}

	// Проверяем, не зарегистрирован ли уже пользователь
	for _, attendee := range attendees {
		if attendee == userEmail {
			return nil, fmt.Errorf("user already registered for this event")
		}
	}

	// Добавляем пользователя в список участников
	attendees = append(attendees, userEmail)

	// Обновляем список участников в базе данных
	_, err = repo.DB.Exec("UPDATE events SET attendees=$1 WHERE id=$2", pq.Array(attendees), eventID)
	if err != nil {
		return nil, err
	}

	return &event, nil
}
