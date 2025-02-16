package utils

import (
	"bytes"
	"html/template"
	"log"
	"os"

	gomail "gopkg.in/gomail.v2"
)

type EventInfo struct {
	Name        string
	Title       string
	Date        string
	Time        string
	Venue       string
	Description string
}

// SendEventRegistrationMail отправляет письмо с подтверждением регистрации на событие
func (i EventInfo) SendEventRegistrationMail() {
	dir, _ := os.Getwd()
	log.Println("Current directory:", dir)

	// Загружаем шаблон письма
	t, err := template.ParseFiles(dir + "/internal/utils/event_template.html")
	if err != nil {
		log.Fatalf("Ошибка при загрузке шаблона: %v", err)
	}

	// Подготавливаем данные для шаблона
	data := struct {
		Name        string
		Title       string
		Date        string
		Time        string
		Venue       string
		Description string
	}{
		Name:        i.Name,
		Title:       i.Title,
		Date:        i.Date,
		Time:        i.Time,
		Venue:       i.Venue,
		Description: i.Description,
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		log.Fatalf("Ошибка при выполнении шаблона: %v", err)
	}

	result := tpl.String()
	m := gomail.NewMessage()
	m.SetHeader("From", "karakuzov.bekbolat@mail.ru")
	m.SetHeader("To", i.Name)
	m.SetHeader("Subject", "Event Registration Confirmation")
	m.SetBody("text/html", result)

	d := gomail.NewDialer("smtp.mail.ru", 465, "karakuzov.bekbolat@mail.ru", "mvarJHYCgX7WcHJ44jKT")

	if err := d.DialAndSend(m); err != nil {
		log.Fatalf("Ошибка при отправке письма: %v", err)
	}
}
