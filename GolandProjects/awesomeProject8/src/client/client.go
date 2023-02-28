package main

import (
	"awesomeProject8/src/proto"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/skorobogatov/input"
	_ "image"
	"io/ioutil"
	"net"
)

// interact - функция, содержащая цикл взаимодействия с сервером.
func interact(conn *net.TCPConn) {
	defer conn.Close()
	encoder, decoder := json.NewEncoder(conn), json.NewDecoder(conn)
	for {
		// Чтение команды из стандартного потока ввода
		fmt.Printf("command = ")
		command := input.Gets()

		// Отправка запроса.
		switch command {
		case "loadImage":
			fmt.Print("Image path = ")
			path := input.Gets()
			reader, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Print("Image not found")
				continue
			}
			encoded := base64.StdEncoding.EncodeToString(reader)
			sendRequest(encoder, "loadImage", &proto.Image{encoded})
		case "quit":
			sendRequest(encoder, "quit", nil)
			return
		case "getSize":
			sendRequest(encoder, "getSize", nil)
		case "getColor":
			fmt.Print("x = ")
			var x, y int
			fmt.Scan(&x)
			fmt.Print("y = ")
			fmt.Scan(&y)
			sendRequest(encoder, "getColor", &proto.Coordinates{x, y})
		default:
			fmt.Printf("error: unknown command\n")
			continue
		}

		// Получение ответа.
		var resp proto.Response
		if err := decoder.Decode(&resp); err != nil {
			fmt.Printf("error: %v\n", err)
			break
		}

		// Вывод ответа в стандартный поток вывода.
		switch resp.Status {
		case "ok":
			fmt.Printf("ok\n")
		case "failed":
			if resp.Data == nil {
				fmt.Printf("error: data field is absent in response\n")
			} else {
				var errorMsg string
				if err := json.Unmarshal(*resp.Data, &errorMsg); err != nil {
					fmt.Printf("error: malformed data field in response\n")
				} else {
					fmt.Printf("failed: %s\n", errorMsg)
				}
			}
		case "getSizeResult":
			if resp.Data == nil {
				fmt.Printf("error: data field is absent in response\n")
			} else {
				var size proto.ImageSize
				if err := json.Unmarshal(*resp.Data, &size); err != nil {
					fmt.Printf("error: malformed data field in response\n")
				} else {
					fmt.Printf("result: %dx%d\n", size.Weight, size.Height)
				}
			}
		case "getColorResult":
			if resp.Data == nil {
				fmt.Printf("error: data field is absent in response\n")
			} else {
				var color proto.ImageColor
				if err := json.Unmarshal(*resp.Data, &color); err != nil {
					fmt.Printf("error: malformed data field in response\n")
				} else {
					fmt.Printf("R: %d G: %d B: %d A: %d\n", color.R/255, color.G/255, color.B/255, color.A/255)
				}
			}
		default:
			fmt.Printf("error: server reports unknown status %q\n", resp.Status)
		}
	}
}

// send_request - вспомогательная функция для передачи запроса с указанной командой
// и данными. Данные могут быть пустыми (data == nil).
func sendRequest(encoder *json.Encoder, command string, data interface{}) {
	var raw json.RawMessage
	raw, _ = json.Marshal(data)
	encoder.Encode(&proto.Request{command, &raw})
}

func main() {
	// Работа с командной строкой, в которой может указываться необязательный ключ -addr.
	var addrStr string
	flag.StringVar(&addrStr, "addr", "127.0.0.1:6000", "specify ip address and port")
	flag.Parse()

	// Разбор адреса, установка соединения с сервером и
	// запуск цикла взаимодействия с сервером.
	if addr, err := net.ResolveTCPAddr("tcp", addrStr); err != nil {
		fmt.Printf("error: %v\n", err)
	} else if conn, err := net.DialTCP("tcp", nil, addr); err != nil {
		fmt.Printf("error: %v\n", err)
	} else {
		interact(conn)
	}
}
