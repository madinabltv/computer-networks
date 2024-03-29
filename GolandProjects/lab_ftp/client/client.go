package main

import (
	"bytes"
	"fmt"
	"github.com/jlaffaye/ftp"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	//fmt.Println("dgfdgfd")
	// постучались
	connection, err := ftp.Dial("students.yss.su:21")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected")

	// авторизация
	err = connection.Login("ftpiu8", "3Ru7yOTA")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Login")

	// создание папки
	err = connection.MakeDir("/BaltaevaMadina")
	if err != nil {
		log.Fatal(err)
	}
	// переход в папку
	connection.ChangeDir("/BaltaevaMadina")

	err = connection.MakeDir("/BaltaevaMadina/dir1")
	if err != nil {
		log.Fatal(err)
	}
	// переход в папку
	connection.ChangeDir("/BaltaevaMadina/dir1")

	// создание буферной строки и запись ее в файл
	data := bytes.NewBufferString("I done it!")
	// загрузка файла на сервер
	connection.Stor("BaltaevaMadina.txt", data)
	// переход в другую папку
	connection.ChangeDir("/BaltaevaMadina")
	// создание буферной строки и запись ее в файл
	data = bytes.NewBufferString("Download content")
	// загрузка файла на сервер
	connection.Stor("news.txt", data)

	// открытие файла
	r, err := connection.Retr("news.txt")
	if err != nil {
		log.Fatal(err)
	}

	// чтение из файла
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	// закрытие файла
	r.Close()

	// создание файла
	file, err := os.Create("news.txt")
	if err != nil {
		log.Fatal(err)
	}
	// запись прочитанного в файл
	file.Write(buf)
	// закрытие файла
	file.Close()

	// построение дерева
	dirLookup("/", connection)

	// удаление папки (рекурсивное)
	connection.RemoveDirRecur("/BaltaevaMadina")
	fmt.Println("Directory deleted")

	// выход
	if err := connection.Quit(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Disconnected")
}

func dirLookup(path string, connection *ftp.ServerConn) {
	fmt.Print("+──" + path)
	dirRec(path, connection, 1)
	fmt.Println()
}

func dirRec(path string, connection *ftp.ServerConn, depth int) {
	entries, err := connection.List(path)
	if err != nil {
		fmt.Println("error in func:", err)
		return
	}
	for _, a := range entries {
		fmt.Println()
		fmt.Print("+")
		for i := 0; i < depth; i++ {
			fmt.Print("   ")
		}
		fmt.Print("└──" + a.Name)
		if strings.Compare(a.Type.String(), "folder") == 0 && a.Name != "." && a.Name != ".." {
			dirRec(path+"/"+a.Name, connection, depth+1)
		}
	}
}
