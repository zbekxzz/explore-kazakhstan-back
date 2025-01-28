package utils

import (
	"database/sql"
	"fmt"
	"github.com/xuri/excelize/v2"
	"strings"
	"time"
)

func ProcessExcelFile(filePath string, db *sql.DB) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	tableName := time.Now().Format("2006_01_02")
	escapedTableName := fmt.Sprintf(`"%s"`, tableName)

	sheetName := f.GetSheetName(0) // Имя первого листа
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("failed to read rows: %w", err)
	}

	if len(rows) == 0 {
		return fmt.Errorf("Excel file is empty")
	}

	// Проверяем заголовки колонок
	expectedHeaders := []string{"time", "complain_type", "complain_description"}
	if len(rows[0]) < len(expectedHeaders) {
		return fmt.Errorf("missing required columns in header")
	}

	for i, header := range expectedHeaders {
		if strings.TrimSpace(rows[0][i]) != header {
			return fmt.Errorf("invalid header at column %d: expected '%s', got '%s'", i+1, header, rows[0][i])
		}
	}

	// Создаем таблицу в базе данных
	createTableQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			time TIMESTAMP,
			complain_type TEXT,
			complain_description TEXT
		)
	`, escapedTableName)
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	if err != nil {
		return fmt.Errorf("failed to read rows: %w", err)
	}

	// Вставляем данные в таблицу
	for i, row := range rows {
		if i == 0 { // Пропускаем заголовки
			continue
		}

		if len(row) < 3 {
			continue
		}

		// Преобразуем строку времени в формат PostgreSQL
		timeString := row[0]

		// Определяем формат
		var parsedTime time.Time
		if strings.Contains(timeString, "/") { // Если формат содержит слэши
			parsedTime, err = time.Parse("1/2/06 15:04", timeString) // Формат MM/DD/YY HH:MM
		} else {
			parsedTime, err = time.Parse("02.01.2006 15:04:05", timeString) // Формат DD.MM.YYYY HH:MM:SS
		}

		if err != nil {
			return fmt.Errorf("failed to parse time '%s': %w", timeString, err)
		}

		complainType := row[1]
		complainDescription := row[2]

		insertQuery := fmt.Sprintf(
			"INSERT INTO %s (time, complain_type, complain_description) VALUES ($1, $2, $3)",
			escapedTableName,
		)
		_, err = db.Exec(insertQuery, parsedTime, complainType, complainDescription)
		if err != nil {
			return fmt.Errorf("failed to insert data: %w", err)
		}
	}

	return nil
}
