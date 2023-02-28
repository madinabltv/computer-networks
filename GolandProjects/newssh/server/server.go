package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	ssh.Handle(func(s ssh.Session) {
		term := terminal.NewTerminal(s, "> ")
		var line string
		for {
			line, _ = term.ReadLine()
			fmt.Println(line)
			str := strings.Split(line, " ")
			cmd := exec.Command(str[0], str[1:]...)
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()

			if err != nil {
				io.WriteString(s, fmt.Sprintf("%s\n", err))
			} else {
				io.WriteString(s, out.String())
			}
		}
	})

	data := map[string]string{"madina": "qwerty123"}

	log.Fatal(ssh.ListenAndServe(":22", nil,
		ssh.PasswordAuth(func(ctx ssh.Context, pass string) bool {
			return pass != "" && data[ctx.User()] == pass
		}),
	))
}

// mkdir <директория>				создание директории на удаленном SSH-сервере;
// rmdir <директория>				удаление директории на удаленном SSH-сервере;
// dir <директория>					вывод содержимого директории;
// mv <путь-к-файлу> <директория>	перемещение файлов из одной директории в другую;
// rm <путь-к-файлу>				удаление файла по имени;
// <путь-к-приложению>				вызов внешних приложений, например ping.
