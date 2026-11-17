package advertisements

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/utils"
)

// GET-запрос поиска объявлений по вхождению текста в название
func HttpGetAdvertisements(t *testing.T, token, textForSearch string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Expect(t).
		Status(http.StatusOK).
		End().Response
}

// POST-запрос создания объявления с фотографиями
func HttpPostAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) *http.Response {
	body, contentType := buildCreateAdvertisementForm(t, request)

	return apiRunner.GetRunner().Auth(token).Create().
		Post(path.AdvertisementPath).
		ContentType(contentType).
		Body(body).
		Expect(t).
		End().Response
}

// GET-запрос получения объявления по ID
func HttpGetAdvertisementByID(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}

// DELETE-запрос удаления объявления по ID
func HttpDeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}

// GET-запрос получения списка фотографий объявления
func HttpGetAdvertisementPhotos(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(fmt.Sprintf(path.AdvertisementPhotosPath, id)).
		Expect(t).
		End().Response
}

// buildCreateAdvertisementForm собирает тело multipart-запроса и возвращает его вместе с Content-Type (с boundary)
func buildCreateAdvertisementForm(t *testing.T, request models.CreateAdvertisementRequest) (string, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fields := map[string]string{
		"title":       request.Title,
		"description": request.Description,
		"price":       strconv.FormatFloat(request.Price, 'f', 2, 64),
		"quantity":    strconv.Itoa(request.Quantity),
	}
	for name, value := range fields {
		require.NoErrorf(t, writer.WriteField(name, value), "Не удалось записать поле формы %s", name)
	}

	for i := range request.PhotosCount {
		header := textproto.MIMEHeader{}
		header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="photos"; filename="photo%d.png"`, i+1))
		header.Set("Content-Type", "image/png")

		part, err := writer.CreatePart(header)
		require.NoError(t, err, "Не удалось создать часть формы для фотографии")
		_, err = part.Write(utils.GeneratePNG())
		require.NoError(t, err, "Не удалось записать фотографию в форму")
	}

	require.NoError(t, writer.Close(), "Не удалось закрыть multipart-форму")

	return body.String(), writer.FormDataContentType()
}
