package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	Config *ssh.ClientConfig
	Host   string
	Port   int
}

func (client *SSHClient) newSession() (*ssh.Session, error) {
	connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", client.Host, client.Port), client.Config)
	if err != nil {
		return nil, fmt.Errorf("Failed to dial: %s", err)
	}

	session, err := connection.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Failed to create session: %s", err)
	}

	modes := ssh.TerminalModes{
		// ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		return nil, fmt.Errorf("request for pseudo terminal failed: %s", err)
	}

	return session, nil
}

func main() {
	sshConfig := &ssh.ClientConfig{
		User: "test",
		Auth: []ssh.AuthMethod{
			ssh.Password("SDHBCXdsedfs222"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client := &SSHClient{
		Config: sshConfig,
		Host:   "151.248.113.144",
		Port:   443,
	}

	var (
		cmd     string
		session *ssh.Session
		err     error
	)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(sshConfig.User, "@izobretarium:~$ ")
		cmd, _ = reader.ReadString('\n')
		cmd = cmd[:len(cmd)-2]

		if session, err = client.newSession(); err != nil {
			fmt.Fprintf(os.Stderr, "session error: %s\n", err)
			os.Exit(1)
		}
		defer session.Close()

		stdout, err := session.StdoutPipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to setup stdout for session: %v", err)
			os.Exit(1)
		}
		stdout_ready := make(chan int)
		go func() {
			io.Copy(os.Stdout, stdout)
			stdout_ready <- 1
		}()

		stderr, err := session.StderrPipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to setup stderr for session: %v", err)
			os.Exit(1)
		}
		stderr_ready := make(chan int)
		go func() {
			io.Copy(os.Stderr, stderr)
			stderr_ready <- 1
		}()

		session.Run(cmd)

		<-stdout_ready
		<-stderr_ready
	}
}

// mkdir <директория>				создание директории на удаленном SSH-сервере;
// rmdir <директория>				удаление директории на удаленном SSH-сервере;
// dir <директория>					вывод содержимого директории;
// mv <путь-к-файлу> <директория>	перемещение файлов из одной директории в другую;
// rm <путь-к-файлу>				удаление файла по имени;
// <путь-к-приложению>				вызов внешних приложений, например ping.