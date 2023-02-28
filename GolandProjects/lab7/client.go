package main

import (
	"crypto/tls"
	"database/sql"
	"encoding/base64"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/smtp"
	"os"
	"strings"
)

func main() {
	var (
		serverAddr = "smtp.gmail.com"
		password   = "galgwqfvbwqgpvew"
		emailAddr  = "madiqwerty2003@gmail.com"
		portNumber = 465
		tos        = []string{
			"blablabla@erunda.com",
			"madina_bltv@vk.com",
		}
		cc                 []string
		attachmentFilePath = "madina.txt"
		filename           = "madina2.txt"
		delimeter          = "**=myohmy689407924327"
		subject            = "Baltaeva_Madina_lab7"
		name               = "Данила Павлович"
		text               = "i love golang <3"
	)
	tlsConfig := tls.Config{
		ServerName:         serverAddr,
		InsecureSkipVerify: true,
	}
	log.Println("Establish TLS connection")
	conn, connErr := tls.Dial("tcp", fmt.Sprintf("%s:%d", serverAddr, portNumber), &tlsConfig)
	if connErr != nil {
		log.Panic(connErr)
	}
	defer conn.Close()
	log.Println("create new email client")
	client, clientErr := smtp.NewClient(conn, serverAddr)
	if clientErr != nil {
		log.Panic(clientErr)
	}
	defer client.Close()
	log.Println("setup authenticate credential")
	auth := smtp.PlainAuth("", emailAddr, password, serverAddr)
	if err := client.Auth(auth); err != nil {
		log.Panic(err)
	}
	log.Println("Start write mail content")
	log.Println("Set 'FROM'")
	if err := client.Mail(emailAddr); err != nil {
		log.Panic(err)
	}
	log.Println("Set 'TO(s)'")
	for _, to := range tos {
		if err := client.Rcpt(to); err != nil {
			log.Panic(err)
		}
	}
	writer, writerErr := client.Data()
	if writerErr != nil {
		log.Panic(writerErr)
	}
	//basic email headers
	sampleMsg := fmt.Sprintf("From: %s\r\n", emailAddr)
	sampleMsg += fmt.Sprintf("To: %s\r\n", strings.Join(tos, ";"))
	if len(cc) > 0 {
		sampleMsg += fmt.Sprintf("Cc: %s\r\n", strings.Join(cc, ";"))
	}
	sampleMsg += "Subject: " + subject + "\r\n"
	log.Println("Mark content to accept multiple contents")
	sampleMsg += "MIME-Version: 1.0\r\n"
	sampleMsg += fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n",
		delimeter)
	//place HTML message
	log.Println("Put HTML message")
	sampleMsg += fmt.Sprintf("\r\n--%s\r\n", delimeter)
	sampleMsg += "Content-Type: text/html; charset=\"utf-8\"\r\n"
	sampleMsg += "Content-Transfer-Encoding: 7bit\r\n"

	sampleMsg += fmt.Sprintf("\r\n%s", "<html>"+
		"<body style=\"background:#CCCCFF\">"+
		"<h1 style=\"background:#FFCCCC\">Hi "+name+"</h1>"+
		"<i style=\"background:#FFCCCC\">"+text+"</i></body></html>\r\n")
	//place file
	log.Println("Put file attachment")
	sampleMsg += fmt.Sprintf("\r\n--%s\r\n", delimeter)
	sampleMsg += "Content-Type: text/plain; charset=\"utf-8\"\r\n"
	sampleMsg += "Content-Transfer-Encoding: base64\r\n"
	sampleMsg += "Content-Disposition: attachment;filename=\"" + filename +
		"\"\r\n"
	//read file
	rawFile, fileErr := os.ReadFile(attachmentFilePath)
	if fileErr != nil {
		log.Panic(fileErr)
	}
	sampleMsg += "\r\n" + base64.StdEncoding.EncodeToString(rawFile)
	//write into email client stream writter
	log.Println("Write content into client writter I/O")
	if _, err := writer.Write([]byte(sampleMsg)); err != nil {
		log.Panic(err)
	}
	if closeErr := writer.Close(); closeErr != nil {
		log.Panic(closeErr)
	}
	err := client.Quit()
	if err != nil {
		log.Panic(err)
	}
	log.Print("done.")
	db, err := sql.Open("mysql", "iu9networkslabs"+":"+"Je2dTYr6"+"@tcp("+"students.yss.su"+")/"+"iu9networkslabs")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err1 := db.Exec("insert into iu9networkslabs.iu9madinabltvSMTP"+" (name, email, msg)"+
		" values (?, ?, ?)",
		name, tos[0], text)
	if err1 != nil {
		panic(err1)
	}
}
