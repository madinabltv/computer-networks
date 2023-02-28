package main

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"github.com/mmcdole/gofeed"
	"log"
	"os"
	"time"
)

const (
	URL  string = "https://news.rambler.ru/rss/Guadeloupe/"
	host string = "students.yss.su"
	port int    = 21
	user string = "ftpiu8"
	pass string = "3Ru7yOTA"
)

func main() {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(URL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Parsed successfully")

	file, err := os.Create("news.txt")
	if err != nil {
		log.Fatal(err)
	}

	data := ""
	for _, item := range feed.Items {
		data += item.Title + ", " + item.Published + ", " + item.Link + "\n"
	}
	file.Write([]byte(data))
	err = file.Close()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("File is written successfully")

	file, err = os.Open("news.txt")
	c, err := ftp.Dial(fmt.Sprintf("%s:%d", host, port), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	err = c.Login(user, pass)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		fmt.Println(err)
	}

	err = c.Stor(fmt.Sprintf("%s_%s_%s.txt", "Baltaeva", "Madina", time.Now().Format("02_Jan_2006_15:04:05")), file)
	if err != nil {
		fmt.Println(err)
	}
	err = file.Close()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("File is uploaded successfully")
}
