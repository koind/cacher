package handler

import (
	"github.com/gin-gonic/gin"
)

// Структура ответа
type response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

// Метод для генерации неуспешного ответа
func responseError(message string) gin.H {
	return gin.H{"status": "error", "message": message}
}

// Метод для генерации успешного ответа
func responseSuccess(data interface{}) response {
	return response{"ok", data}
}
