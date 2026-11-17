package advertisements

import (
	"fmt"
	"testing"

	httpClient "api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/utils"
)

type createdAdvertisement struct {
	token string
	id    string
}

// Fixture создает объявления для тестов и запоминает их, чтобы удалить после теста (тесты остаются идемпотентными)
type Fixture struct {
	created []createdAdvertisement
}

// CreateAdvertisement создает объявление от имени владельца токена и ставит его на удаление в Cleanup
func (f *Fixture) CreateAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) models.Advertisement {
	ad := CreateAdvertisement(t, token, request)
	f.created = append(f.created, createdAdvertisement{token: token, id: ad.ID})
	utils.LogWithLabelAndTimeStamp("Fixture", fmt.Sprintf("Создано объявление id=%s, title=%s", ad.ID, ad.Title))

	return ad
}

// Cleanup удаляет все созданные через фикстуру объявления.
func (f *Fixture) Cleanup(t *testing.T) {
	for _, ad := range f.created {
		utils.LogWithLabelAndTimeStamp("Cleanup", fmt.Sprintf("Удаляем объявление id=%s", ad.id))
		raw := httpClient.HttpDeleteAdvertisement(t, ad.token, ad.id)
		_ = raw.Body.Close()
	}
	f.created = nil
}
