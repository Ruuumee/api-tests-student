package advertisements

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	httpClient "api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/utils"
)

//константы для создания нового обьявляния
const (
	DefaultDescription = "Объявление про что-то и то-то"
	DefaultPrice       = 1730.29
	DefaultQuantity    = 3
	DefaultPhotosCount = 1
	titlePrefix        = "AutotestAvito_QA-"
	titleRandomLength  = 10
)

// NewAdvertisementRequest возвращает корректный запрос на создание объявления со всеми параметрами и уникальным названием
func NewAdvertisementRequest() models.CreateAdvertisementRequest {
	return models.CreateAdvertisementRequest{
		Title:       titlePrefix + utils.RandomString(titleRandomLength),
		Description: DefaultDescription,
		Price:       DefaultPrice,
		Quantity:    DefaultQuantity,
		PhotosCount: DefaultPhotosCount,
	}
}

// GetAdvertisementsList ищет объявления по вхождению текста в название
func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	raw := httpClient.HttpGetAdvertisements(t, token, text)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	return response.ParseJSON[models.AdvertisementsList](t, body)
}

// CreateAdvertisement создает объявление и возвращает его, ожидает 201
func CreateAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) models.Advertisement {
	raw := httpClient.HttpPostAdvertisement(t, token, request)
	body := response.ReadBody(t, raw, http.StatusCreated)

	ad := response.ParseJSON[models.Advertisement](t, body)
	require.NotEmptyf(t, ad.ID, "В ответе на создание объявления нет id: %s", string(body))

	return ad
}

// CreateAdvertisementExpectError пытается создать объявление
func CreateAdvertisementExpectError(t *testing.T, token string, request models.CreateAdvertisementRequest,
	expectedStatusCode int) response.ErrorResponse {
	raw := httpClient.HttpPostAdvertisement(t, token, request)
	body := response.ReadBody(t, raw, expectedStatusCode)

	return response.ParseJSON[response.ErrorResponse](t, body)
}

// GetAdvertisementByID возвращает объявление по ID, ожидает 200
func GetAdvertisementByID(t *testing.T, token, id string) models.Advertisement {
	raw := httpClient.HttpGetAdvertisementByID(t, token, id)
	body := response.ReadBody(t, raw, http.StatusOK)

	return response.ParseJSON[models.Advertisement](t, body)
}

// GetAdvertisementPhotos возвращает фотографии объявления, ожидает 200
func GetAdvertisementPhotos(t *testing.T, token, id string) []models.Photo {
	raw := httpClient.HttpGetAdvertisementPhotos(t, token, id)
	body := response.ReadBody(t, raw, http.StatusOK)

	return response.ParseJSON[[]models.Photo](t, body)
}

// DeleteAdvertisement удаляет объявление по ID, ожидает 204
func DeleteAdvertisement(t *testing.T, token, id string) {
	raw := httpClient.HttpDeleteAdvertisement(t, token, id)
	response.ReadBody(t, raw, http.StatusNoContent)
}
