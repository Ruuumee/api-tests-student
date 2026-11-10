package advertisements

import (
	"encoding/json"
	"io"
	"testing"

	httpClient "api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"github.com/stretchr/testify/require"
)

func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	
	// TODO 1: вызвать advertisements.HttpGetAdvertisements
	raw := httpClient.HttpGetAdvertisements(t, token, text)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	// TODO 2: распарсить body в models.AdvertisementsList через json.Unmarshal
	resp := models.AdvertisementsList{}
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}