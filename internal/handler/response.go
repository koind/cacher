package handler

// Структура ответа
type response struct {
	// Статус ответа
	Status string `json:"status"`
	// Тело ответа
	Data interface{} `json:"data,omitempty"`
	// Текста ошибки
	Message string `json:"message,omitempty"`
}

// Метод для генерации неуспешного ответа
func responseError(err error) response {
	return response{Status: "error", Message: err.Error()}
}

// Метод для генерации успешного ответа
func responseSuccess(data interface{}) response {
	return response{Status: "ok", Data: data}
}
