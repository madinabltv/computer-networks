package main

import (
	"fmt"
	"log"
	"net/smtp"
)

func main() {

	//mail info
	from := Email
	password := MailPassword

	//smtp info
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	//receiver
	to := []string{"danila@posevin.com", "madina_bltv@vk.com"}

	///message

	subject := "Subject:Baltaeva_Madina_2022-10-25_17:23\n"

	body := "i love golang"
	message := []byte(subject + body)

	//authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	//send mail

	//smtp.gmail.com:587
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("mail send")
}
