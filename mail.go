package main

import (
	"fmt"
	"net/smtp"
)

func sendMail(email string) {

	from := "centre-vaccination@gmail.com"
	password := "passwordsupersecure"

	to := []string{
		email,
	}

	// configuration smtp.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte("Votre rendez vous est validé")

	// Authentification.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Envoi email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email envoyé!")
}
