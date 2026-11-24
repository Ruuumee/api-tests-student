package advertisements

import (
	"bytes"
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisements/models"
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

	price := strconv.FormatFloat(request.Price, 'f', 2, 64)
	if request.PriceRaw != "" {
		price = request.PriceRaw
	}
	quantity := strconv.Itoa(request.Quantity)
	if request.QuantityRaw != "" {
		quantity = request.QuantityRaw
	}

	fields := map[string]string{
		"title":       request.Title,
		"description": request.Description,
		"price":       price,
		"quantity":    quantity,
	}
	for name, value := range fields {
		require.NoErrorf(t, writer.WriteField(name, value), "Не удалось записать поле формы %s", name)
	}

	for _, photoPath := range request.PhotoPaths {
		content, err := os.ReadFile(filepath.Clean(photoPath))
		require.NoErrorf(t, err, "Не удалось прочитать файл фотографии %s", photoPath)

		contentType := mime.TypeByExtension(filepath.Ext(photoPath))
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		header := textproto.MIMEHeader{}
		header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="photos"; filename="%s"`, filepath.Base(photoPath)))
		header.Set("Content-Type", contentType)

		part, err := writer.CreatePart(header)
		require.NoErrorf(t, err, "Не удалось создать часть формы для фотографии %s", photoPath)
		_, err = part.Write(content)
		require.NoErrorf(t, err, "Не удалось записать фотографию %s в форму", photoPath)
	}

	require.NoError(t, writer.Close(), "Не удалось закрыть multipart-форму")

	return body.String(), writer.FormDataContentType()
}
