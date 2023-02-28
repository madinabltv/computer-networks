package proto

import "encoding/json"

// Request -- запрос клиента к серверу.
type Request struct {
	// Поле Command может принимать три значения:
	// * "quit" - прощание с сервером (после этого сервер рвёт соединение);
	// * "loadImage" - передача изображения;
	// * "getSize" - просьба посчитать размер изображения;
	// * "getColor" - просьба вывести цвет пикселя;
	Command string `json:"command"`

	// Если Command == "add", в поле Data должна лежать дробь
	// в виде структуры Fraction.
	// В противном случае, поле Data пустое.
	Data *json.RawMessage `json:"data"`
}

// Response -- ответ сервера клиенту.
type Response struct {
	// Поле Status может принимать три значения:
	// * "ok" - успешное выполнение команды "quit" или "add";
	// * "failed" - в процессе выполнения команды произошла ошибка;
	// * "result" - среднее арифметическое дробей вычислено.
	Status string `json:"status"`

	// Если Status == "failed", то в поле Data находится сообщение об ошибке.
	// Если Status == "result", в поле Data должна лежать дробь
	// в виде структуры Fraction.
	// В противном случае, поле Data пустое.
	Data *json.RawMessage `json:"data"`
}
type Image struct {
	Encoded string `json:"encoded"`
}
type Coordinates struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type ImageColor struct {
	//Request string `json:"request"`
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
	A int `json:"a"`
}
type ImageSize struct {
	//Request string `json:"request"`
	Weight int `json:"weight"`
	Height int `json:"height"`
}
