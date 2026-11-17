package response

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// ErrorResponse описывает тело ответа с ошибкой (поля error и message из swagger)
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// ReadBody проверяет статус-код ответа и возвращает тело ответа.
// Если статус не совпал, тело ответа попадает в сообщение об ошибке, чтобы было понятно, почему упал тест
func ReadBody(t *testing.T, resp *http.Response, expectedStatusCode int) []byte {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Не удалось прочитать тело ответа")
	require.Equalf(t, expectedStatusCode, resp.StatusCode,
		"Неожиданный HTTP status code, тело ответа: %s", string(body))

	return body
}

// ParseJSON разбирает тело ответа в структуру указанного типа
func ParseJSON[T any](t *testing.T, body []byte) T {
	var result T
	require.NoErrorf(t, json.Unmarshal(body, &result), "Не удалось разобрать JSON: %s", string(body))

	return result
}
