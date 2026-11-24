package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/auth"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
)

const (
	randomUserFirstName = "Autotest"
	randomPartLength    = 10
)

// Login выполняет авторизацию пользователя и возвращает токен с данными пользователя, ожидает 200
func Login(t *testing.T, username, password string) models.LoginOkResponse {
	require.NotEmpty(t, username, "Отсутствует логин пользователя")
	require.NotEmpty(t, password, "Отсутствует пароль пользователя")

	request := models.LoginRequest{
		Email:    username,
		Password: password,
	}

	loginRawResponse := auth.HttpPostAuthLogin(t, request)
	defer loginRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(loginRawResponse.Body)
	require.NoError(t, errResponse, "Не удалось прочитать тело ответа на авторизацию")

	loginResponse := models.LoginOkResponse{}
	err := json.Unmarshal(respArray, &loginResponse)
	require.NoErrorf(t, err, "Не удалось разобрать ответ на авторизацию: %s", string(respArray))

	return loginResponse
}

// Register регистрирует нового пользователя и возвращает токен с данными пользователя, ожидает 201
func Register(t *testing.T, request models.RegisterRequest) models.LoginOkResponse {
	raw := auth.HttpPostAuthRegister(t, request)
	body := response.ReadBody(t, raw, http.StatusCreated)

	registerResponse := response.ParseJSON[models.LoginOkResponse](t, body)
	require.NotEmptyf(t, registerResponse.Token, "В ответе на регистрацию нет токена: %s", string(body))

	return registerResponse
}

// RegisterRandomUser регистрирует нового пользователя со случайными email и паролем
func RegisterRandomUser(t *testing.T) models.LoginOkResponse {
	request := models.RegisterRequest{
		Email:     strings.ToLower(fmt.Sprintf("autotest_%s@example.com", utils.RandomString(randomPartLength))),
		FirstName: randomUserFirstName,
		Password:  "Aa1_" + utils.RandomString(randomPartLength),
	}

	return Register(t, request)
}
