package main

import (
	"fmt"
	"net/smtp"
	"os"
)

func main() {
	message := []byte(fmt.Sprintf("Subject: %s\n%s, %s", os.Args[2], os.Args[4], os.Args[3]))
	auth := smtp.PlainAuth("", "madiqwerty2003@gmail.com", "galgwqfvbwqgpvew", "smtp.gmail.ru")
	err := smtp.SendMail("smtp.gmail.ru:587", auth, "madiqwerty2003@gmail.com", []string{os.Args[1]}, message)
	if err != nil {
		panic(err)
	}
}
