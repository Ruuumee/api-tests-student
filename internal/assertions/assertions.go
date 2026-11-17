package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/advertisements/models"
	favoriteModels "api-tests-template/internal/managers/favorites/models"
)

const (
	UnauthorizedErrorCode = "unauthorized"
	priceDelta            = 0.001
)

func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	for _, item := range items {
		if item.Title == title {
			return
		}
	}
	require.Fail(t, "Объявление не найдено в списке", "Искали заголовок: '%s', но его нет в ответе", title)
}

func AdvertisementMatchesRequest(t *testing.T, ad models.Advertisement, request models.CreateAdvertisementRequest, ownerID string) {
	require.Equal(t, request.Title, ad.Title, "title не совпадает с переданным при создании")
	require.Equal(t, request.Description, ad.Description, "description не совпадает с переданным при создании")
	require.InDelta(t, request.Price, ad.Price, priceDelta, "price не совпадает с переданным при создании")
	require.Equal(t, request.Quantity, ad.Quantity, "quantity не совпадает с переданным при создании")
	require.Equal(t, ownerID, ad.UserID, "user_id не совпадает с пользователем, создавшим объявление")
}

// в ответе с ошибкой заполнены поля error и message?
func AssertErrorResponse(t *testing.T, resp response.ErrorResponse) {
	require.NotEmpty(t, resp.Error, "В ответе с ошибкой не заполнено поле error")
	require.NotEmpty(t, resp.Message, "В ответе с ошибкой не заполнено поле message")
}

// проверяет тело ответа 401: заполнены error и message, error = unauthorized
func AssertUnauthorizedResponse(t *testing.T, resp response.ErrorResponse) {
	AssertErrorResponse(t, resp)
	require.Equal(t, UnauthorizedErrorCode, resp.Error, "Некорректное значение error параметра")
}

// объявление есть в списке избранного
func FavoritesContainAdvertisement(t *testing.T, items []favoriteModels.Favorite, adID string) {
	require.Truef(t, hasFavorite(items, adID),
		"Объявление %s должно быть в избранном, но его нет в ответе (всего в списке: %d)", adID, len(items))
}

// объявления нет в списке избранного
func FavoritesNotContainAdvertisement(t *testing.T, items []favoriteModels.Favorite, adID string) {
	require.Falsef(t, hasFavorite(items, adID),
		"Объявление %s не должно быть в избранном, но оно есть в ответе", adID)
}

func hasFavorite(items []favoriteModels.Favorite, adID string) bool {
	for _, item := range items {
		if item.Advertisement.ID == adID {
			return true
		}
	}
	return false
}
