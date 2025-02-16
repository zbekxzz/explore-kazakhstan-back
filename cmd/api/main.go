package main

import (
	"auth/config"
	"auth/internal/database"
	"auth/internal/handlers"
	"auth/internal/jwt"
	"auth/internal/repository"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	cfg := config.AppConfig

	db := database.InitDB(cfg.DatabaseURL)
	userRepo := repository.NewUserRepository(db)
	eventRepo := repository.NewEventRepository(db)

	jwt.InitJWT(cfg.SecretKey, cfg.AccessTokenExpiresTime, cfg.RefreshTokenExpiresTime)

	http.HandleFunc("/api/login", handlers.LoginHandler(userRepo))
	http.HandleFunc("/api/register", handlers.RegisterHandler(userRepo))
	
	http.HandleFunc("/user/list", handlers.ListUsersHandler(userRepo))
	http.HandleFunc("/user/update", handlers.UpdateUserHandler(userRepo))
	http.HandleFunc("/user/delete", handlers.DeleteUserHandler(userRepo))

	http.HandleFunc("/events/register", handlers.RegisterForEventHandler(eventRepo))
	http.HandleFunc("/events/list", handlers.ListEventsHandler(eventRepo))
	http.HandleFunc("/events/get", handlers.GetEventHandler(eventRepo))
	http.HandleFunc("/events/create", handlers.CreateEventHandler(eventRepo))
	http.HandleFunc("/events/update", handlers.UpdateEventHandler(eventRepo))
	http.HandleFunc("/events/delete", handlers.DeleteEventHandler(eventRepo))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173/app/user-directory", "http://localhost:5173"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Authorization",
			"Content-Type",
		},
		AllowCredentials: true,
	})

	handler := c.Handler(http.DefaultServeMux)

	log.Println("Starting server on :3000")
	log.Fatal(http.ListenAndServe(":"+cfg.Port, handler))
}
