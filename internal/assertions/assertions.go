package assertions

import (
	"testing"

	"api-tests-template/internal/managers/advertisements/models"

	"github.com/stretchr/testify/require"
)

func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	for _, item := range items {
		if item.Title == title {
			return
		}
	}
	require.Fail(t, "Объявление не найдено в списке", "Искали заголовок: '%s', но его нет в ответе", title)
}
