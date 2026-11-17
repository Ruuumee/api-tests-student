package favorites

import (
	"net/http"
	"testing"

	httpClient "api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/favorites/models"
)

// AddToFavorites добавляет объявление в избранное и возвращает созданную запись, ожидает 201
func AddToFavorites(t *testing.T, token, adID string) models.Favorite {
	raw := httpClient.HttpPostFavorite(t, token, models.AddFavoriteRequest{AdvertisementID: adID})
	body := response.ReadBody(t, raw, http.StatusCreated)

	return response.ParseJSON[models.Favorite](t, body)
}

// AddToFavoritesExpectError пытается добавить объявление в избранное
func AddToFavoritesExpectError(t *testing.T, token, adID string, expectedStatusCode int) response.ErrorResponse {
	raw := httpClient.HttpPostFavorite(t, token, models.AddFavoriteRequest{AdvertisementID: adID})
	body := response.ReadBody(t, raw, expectedStatusCode)

	return response.ParseJSON[response.ErrorResponse](t, body)
}

// GetFavorites возвращает список избранного пользователя, ожидает 200
func GetFavorites(t *testing.T, token string) models.FavoriteList {
	raw := httpClient.HttpGetFavorites(t, token)
	body := response.ReadBody(t, raw, http.StatusOK)

	return response.ParseJSON[models.FavoriteList](t, body)
}

// GetFavoritesExpectError запрашивает список избранного
func GetFavoritesExpectError(t *testing.T, token string, expectedStatusCode int) response.ErrorResponse {
	raw := httpClient.HttpGetFavorites(t, token)
	body := response.ReadBody(t, raw, expectedStatusCode)

	return response.ParseJSON[response.ErrorResponse](t, body)
}

// RemoveFromFavorites удаляет объявление из избранного, ожидает 204
func RemoveFromFavorites(t *testing.T, token, adID string) {
	raw := httpClient.HttpDeleteFavorite(t, token, adID)
	response.ReadBody(t, raw, http.StatusNoContent)
}

// RemoveFromFavoritesExpectError пытается удалить объявление из избранного
func RemoveFromFavoritesExpectError(t *testing.T, token, adID string, expectedStatusCode int) response.ErrorResponse {
	raw := httpClient.HttpDeleteFavorite(t, token, adID)
	body := response.ReadBody(t, raw, expectedStatusCode)

	return response.ParseJSON[response.ErrorResponse](t, body)
}
