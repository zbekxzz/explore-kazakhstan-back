package handlers

import (
	"auth/internal/repository"
	"auth/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// ListEventsHandler возвращает список всех событий
func ListEventsHandler(repo *repository.EventRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		events, err := repo.GetAllEvents()
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to retrieve events: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(events); err != nil {
			http.Error(w, "failed to encode events", http.StatusInternalServerError)
		}
	}
}

// GetEventHandler возвращает одно событие по ID
func GetEventHandler(repo *repository.EventRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventID := r.URL.Query().Get("id")
		if eventID == "" {
			http.Error(w, "missing event ID", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(eventID)
		if err != nil {
			http.Error(w, "invalid event ID", http.StatusBadRequest)
			return
		}

		event, err := repo.FindByID(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to retrieve event: %v", err), http.StatusInternalServerError)
			return
		}
		if event == nil {
			http.Error(w, "event not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(event)
	}
}

// CreateEventHandler создает новое событие
func CreateEventHandler(repo *repository.EventRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var event repository.Event
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		err := repo.Create(event)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to create event: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "Event created successfully")
	}
}

// UpdateEventHandler обновляет существующее событие
func UpdateEventHandler(repo *repository.EventRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventID := r.URL.Query().Get("id")
		if eventID == "" {
			http.Error(w, "missing event ID", http.StatusBadRequest)
			return
		}

		var event repository.Event
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		err := repo.UpdateByID(eventID, event)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to update event: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Event with ID %s updated successfully", eventID)
	}
}

// DeleteEventHandler удаляет событие по ID
func DeleteEventHandler(repo *repository.EventRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventID := r.URL.Query().Get("id")
		if eventID == "" {
			http.Error(w, "missing event ID", http.StatusBadRequest)
			return
		}

		err := repo.DeleteByID(eventID)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to delete event: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Event with ID %s deleted successfully", eventID)
	}
}

func RegisterForEventHandler(repo *repository.EventRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			EventID string `json:"event_id"`
			Email   string `json:"email"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		event, err := repo.RegisterUserForEvent(request.EventID, request.Email)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to register user: %v", err), http.StatusInternalServerError)
			return
		}

		mailInfo := utils.EventInfo{
			Name:        request.Email,
			Title:       event.Title,
			Date:        event.Date,
			Time:        event.Time,
			Venue:       event.Venue,
			Description: event.Description,
		}
		go mailInfo.SendEventRegistrationMail()

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "User %s registered for event %s", request.Email, request.EventID)
	}
}
