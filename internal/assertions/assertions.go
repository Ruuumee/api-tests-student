package assertions

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/advertisements/models"
	favoriteModels "api-tests-template/internal/managers/favorites/models"
	"api-tests-template/internal/utils"
)

const (
	// UnauthorizedErrorCode значение поля error в ответе 401
	UnauthorizedErrorCode = "unauthorized"
	// MissingAuthMessage значение поля message в ответе 401, если токен не передан
	MissingAuthMessage = "Authorization header required"
	// InvalidAuthMessage значение поля message в ответе 401, если токен неверный
	InvalidAuthMessage = "Invalid or expired token"
	// InvalidRequestErrorCode значение поля error в ответе 400
	InvalidRequestErrorCode = "invalid_request"
	// NotFoundErrorCode значение поля error в ответе 404
	NotFoundErrorCode = "not_found"
	// ForbiddenErrorCode значение поля error в ответе 403
	ForbiddenErrorCode = "forbidden"
	// AlreadyExistsErrorCode значение поля error в ответе 409
	AlreadyExistsErrorCode = "already_exists"
	priceDelta             = 0.001
)

// ContainsAdvertisementByTitle проверяет, что в списке есть объявление с указанным названием
func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	for _, item := range items {
		if item.Title == title {
			return
		}
	}
	require.Fail(t, "Объявление не найдено в списке", "Искали заголовок: '%s', но его нет в ответе", title)
}

// AdvertisementMatchesRequest проверяет, что объявление содержит параметры, переданные при создании
func AdvertisementMatchesRequest(t *testing.T, ad models.Advertisement, request models.CreateAdvertisementRequest, ownerID string) {
	require.Equal(t, request.Title, ad.Title, "title не совпадает с переданным при создании")
	require.Equal(t, request.Description, ad.Description, "description не совпадает с переданным при создании")
	require.InDelta(t, request.Price, ad.Price, priceDelta, "price не совпадает с переданным при создании")
	require.Equal(t, request.Quantity, ad.Quantity, "quantity не совпадает с переданным при создании")
	require.Equal(t, ownerID, ad.UserID, "user_id не совпадает с пользователем, создавшим объявление")
}

// в ответе с ошибкой заполнены поля error и message?
func AssertErrorResponse(t *testing.T, resp response.ErrorResponse) {
	utils.LogWithLabelAndTimeStamp("Response", fmt.Sprintf("error=%q, message=%q", resp.Error, resp.Message))
	require.NotEmpty(t, resp.Error, "В ответе с ошибкой не заполнено поле error")
	require.NotEmpty(t, resp.Message, "В ответе с ошибкой не заполнено поле message")
}

// проверяет, что в ответе с ошибкой error и message равны ожидаемым
func AssertErrorResponseEquals(t *testing.T, resp response.ErrorResponse, wantError, wantMessage string) {
	AssertErrorResponse(t, resp)
	require.Equalf(t, wantError, resp.Error, "Некорректное значение error: ожидали %q, получили %q", wantError, resp.Error)
	require.Equalf(t, wantMessage, resp.Message, "Некорректное значение message: ожидали %q, получили %q", wantMessage, resp.Message)
}

// проверяет тело ответа 401: error = unauthorized, message равен ожидаемому
func AssertUnauthorizedResponse(t *testing.T, resp response.ErrorResponse, wantMessage string) {
	AssertErrorResponseEquals(t, resp, UnauthorizedErrorCode, wantMessage)
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

// hasFavorite возвращает true, если среди избранного есть объявление с указанным ID
func hasFavorite(items []favoriteModels.Favorite, adID string) bool {
	for _, item := range items {
		if item.Advertisement.ID == adID {
			return true
		}
	}
	return false
}
