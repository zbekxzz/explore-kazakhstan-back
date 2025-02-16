package utils

import (
	"bytes"
	"html/template"
	"log"
	"os"

	gomail "gopkg.in/gomail.v2"
)

type Info struct {
	Name string
}

func (i Info) SendMail() {
	// Проверяем текущую директорию
	dir, _ := os.Getwd()
	log.Println("Current directory:", dir)

	t, err := template.ParseFiles(dir + "/internal/utils/template.html")
	if err != nil {
		log.Fatalf("Ошибка при загрузке шаблона: %v", err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, i); err != nil {
		log.Fatalf("Ошибка при выполнении шаблона: %v", err)
	}

	result := tpl.String()
	m := gomail.NewMessage()
	m.SetHeader("From", "karakuzov.bekbolat@mail.ru")
	m.SetHeader("To", "zbekxzz@gmail.com")
	m.SetHeader("Subject", "golang test")
	m.SetBody("text/html", result)

	d := gomail.NewDialer("smtp.mail.ru", 465, "karakuzov.bekbolat@mail.ru", "mvarJHYCgX7WcHJ44jKT")

	if err := d.DialAndSend(m); err != nil {
		log.Fatalf("Ошибка при отправке письма: %v", err)
	}
}
