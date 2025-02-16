package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword() {
	password := "beka"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	fmt.Println(string(hashedPassword))
}
