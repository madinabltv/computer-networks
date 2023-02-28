package main

import (
	"awesomeProject8/src/proto"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mgutz/logxi/v1"
	"image"
	_ "image/jpeg"
	"net"
	"strings"
)

// Client - состояние клиента.
type Client struct {
	logger log.Logger    // Объект для печати логов
	conn   *net.TCPConn  // Объект TCP-соединения
	enc    *json.Encoder // Объект для кодирования и отправки сообщений
	image  image.Image
}

// NewClient - конструктор клиента, принимает в качестве параметра
// объект TCP-соединения.
func NewClient(conn *net.TCPConn) *Client {
	return &Client{
		logger: log.New(fmt.Sprintf("client %s", conn.RemoteAddr().String())),
		conn:   conn,
		enc:    json.NewEncoder(conn),
	}
}

// serve - метод, в котором реализован цикл взаимодействия с клиентом.
// Подразумевается, что метод serve будет вызаваться в отдельной go-программе.
func (client *Client) serve() {
	defer client.conn.Close()
	decoder := json.NewDecoder(client.conn)
	for {
		var req proto.Request
		if err := decoder.Decode(&req); err != nil {
			client.logger.Error("cannot decode message", "reason", err)
			break
		} else {
			client.logger.Info("received command", "command", req.Command)
			if client.handleRequest(&req) {
				client.logger.Info("shutting down connection")
				break
			}
		}
	}
}

// handleRequest - метод обработки запроса от клиента. Он возвращает true,
// если клиент передал команду "quit" и хочет завершить общение.
func (client *Client) handleRequest(req *proto.Request) bool {
	switch req.Command {
	case "quit":
		client.respond("ok", nil)
		return true
	case "loadImage":
		errorMsg := ""
		if req.Data == nil {
			errorMsg = "data field is absent"
		} else {
			var imageSchema proto.Image
			if err := json.Unmarshal(*req.Data, &imageSchema); err != nil {
				errorMsg = "malformed data field"
			} else {
				if image, _, err := image.Decode(base64.NewDecoder(base64.StdEncoding, strings.NewReader(imageSchema.Encoded))); err != nil {
					errorMsg = "malformed data field"
				} else {
					client.logger.Info("Saving image")
					client.image = image
				}
			}
		}
		if errorMsg == "" {
			client.respond("ok", nil)
		} else {
			client.logger.Error("failed", "reason", errorMsg)
			client.respond("failed", errorMsg)
		}
	case "getSize":
		if client.image == nil {
			client.logger.Error("getSize Error", "reason", "image not found")
			client.respond("failed", "No image load")
		} else {
			client.respond("getSizeResult", &proto.ImageSize{
				client.image.Bounds().Max.X,
				client.image.Bounds().Max.Y,
			})
		}
	case "getColor":
		if client.image == nil {
			client.logger.Error("getSize Error", "reason", "image not found")
			client.respond("failed", "No image load")
		} else {
			errorMsg := ""
			if req.Data == nil {
				errorMsg = "data field is absent"
			} else {
				var coordinatesSchema proto.Coordinates
				if err := json.Unmarshal(*req.Data, &coordinatesSchema); err != nil {
					errorMsg = "malformed data field"
				} else {
					r, g, b, a := client.image.At(coordinatesSchema.X, coordinatesSchema.Y).RGBA()
					client.respond("getColorResult", &proto.ImageColor{
						int(r), int(g), int(b), int(a),
					})
				}
			}

			if errorMsg != "" {
				client.logger.Error("failed", "reason", errorMsg)
				client.respond("failed", errorMsg)
			}
		}
	default:
		client.logger.Error("unknown command")
		client.respond("failed", "unknown command")
	}
	return false
}

// respond - вспомогательный метод для передачи ответа с указанным статусом
// и данными. Данные могут быть пустыми (data == nil).
func (client *Client) respond(status string, data interface{}) {
	var raw json.RawMessage
	raw, _ = json.Marshal(data)
	client.enc.Encode(&proto.Response{status, &raw})
}

func main() {
	// Работа с командной строкой, в которой может указываться необязательный ключ -addr.
	var addrStr string
	flag.StringVar(&addrStr, "addr", "127.0.0.1:6000", "specify ip address and port")
	flag.Parse()

	// Разбор адреса, строковое представление которого находится в переменной addrStr.
	if addr, err := net.ResolveTCPAddr("tcp", addrStr); err != nil {
		log.Error("address resolution failed", "address", addrStr)
	} else {
		log.Info("resolved TCP address", "address", addr.String())

		// Инициация слушания сети на заданном адресе.
		if listener, err := net.ListenTCP("tcp", addr); err != nil {
			log.Error("listening failed", "reason", err)
		} else {
			// Цикл приёма входящих соединений.
			for {
				if conn, err := listener.AcceptTCP(); err != nil {
					log.Error("cannot accept connection", "reason", err)
				} else {
					log.Info("accepted connection", "address", conn.RemoteAddr().String())

					// Запуск go-программы для обслуживания клиентов.
					go NewClient(conn).serve()
				}
			}
		}
	}
	// Создаем TCP сервер и слушаем запросы клиентов

}
