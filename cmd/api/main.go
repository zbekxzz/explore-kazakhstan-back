package main

import (
	"auth/config"
	"auth/internal/jwt"
	"auth/internal/repository"
	"github.com/rs/cors"
	"log"
	"net/http"

	"auth/internal/database"
	"auth/internal/handlers"
)

func main() {
	cfg := config.AppConfig

	db := database.InitDB(cfg.DatabaseURL)
	repo := repository.NewUserRepository(db)

	jwt.InitJWT(cfg.SecretKey, cfg.AccessTokenExpiresTime, cfg.RefreshTokenExpiresTime)

	http.HandleFunc("/api/login", handlers.LoginHandler(repo))

	http.HandleFunc("/admin/list", handlers.ListUsersHandler(repo))
	http.HandleFunc("/admin/create", handlers.CreateUserHandler(repo))
	http.HandleFunc("/admin/update", handlers.UpdateUserHandler(repo))
	http.HandleFunc("/admin/delete", handlers.DeleteUserHandler(repo))

	http.HandleFunc("/api/upload", handlers.UploadHandler(db))

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
			"Authorization", // Разрешаем заголовок Authorization
			"Content-Type",  // Разрешаем Content-Type
		},
		AllowCredentials: true,
	})
	handler := c.Handler(http.DefaultServeMux)

	log.Println("Starting server on :3000")
	log.Fatal(http.ListenAndServe(":"+cfg.Port, handler))
}
