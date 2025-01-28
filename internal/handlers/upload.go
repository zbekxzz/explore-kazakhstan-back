package handlers

import (
	"auth/internal/utils"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func UploadHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(10 << 20)

		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		tempDir := "./temp"
		os.MkdirAll(tempDir, os.ModePerm)
		tempFilePath := filepath.Join(tempDir, handler.Filename)

		tempFile, err := os.Create(tempFilePath)
		if err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}
		defer tempFile.Close()

		_, err = tempFile.ReadFrom(file)
		if err != nil {
			http.Error(w, "Failed to write file", http.StatusInternalServerError)
			return
		}

		if err := utils.ProcessExcelFile(tempFilePath, db); err != nil {
			os.Remove(tempFilePath)
			http.Error(w, fmt.Sprintf("Failed to process Excel file: %v", err), http.StatusUnprocessableEntity)
			fmt.Println("Файл не загружен")
			return
		}
		fmt.Println("Файл загружен")
		os.Remove(tempFilePath)

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "File uploaded and processed successfully")
	}
}
