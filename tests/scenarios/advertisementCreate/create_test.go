package advertisementCreate

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	advertisementModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"

	base "api-tests-template/tests"
)

const (
	maxTitleLength = 255
	maxPhotosCount = 3

	serverTitleLimit = 50
	incorrectToken   = "incorrect_token"

	validPhotoPath   = "photos/valid.png"
	invalidPhotoPath = "photos/invalid.txt"

	negativeQuantity = -1
)

// validPhotos возвращает пути к count валидным фотографиям (один и тот же файл прикладывается count раз)
func validPhotos(count int) []string {
	paths := make([]string, count)
	for i := range paths {
		paths[i] = validPhotoPath
	}

	return paths
}

// набор тестов для функционала создания объявлений.
type TestSuite struct {
	suite.Suite
	loginData authModels.LoginOkResponse
	fixture   advertisements.Fixture
}

// точка входа для запуска всех тестов создания объявлений в этом файле.
func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

// Инициализирует окружение, авторизует пользователя
func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователя с кредами из переменных окружения")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

// Выполняется после каждого теста для изоляции даанных
func (s *TestSuite) TearDownTest() {
	s.fixture.Cleanup(s.T())
}

// Выполняет финальную глобальную очистку окружения.
func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// Positive: создаем объявление со всеми параметрами и проверяем GET по ID, поиск по названию и фотографии
func (s *TestSuite) TestCreateAdvertisementPositive() {
	token := s.loginData.Token
	request := advertisements.NewAdvertisementRequest(validPhotos(maxPhotosCount)...)

	var created advertisementModels.Advertisement
	isCreated := s.Run("Создаем объявление со всеми параметрами и фотографией", func() {
		base.Step(fmt.Sprintf("Создаем объявление с названием %s", request.Title))
		created = s.fixture.CreateAdvertisement(s.T(), token, request)

		base.Check("Объявление создано и в ответе есть id")
		require.NotEmptyf(s.T(), created.ID, "В ответе на создание объявления нет id, название: %s", request.Title)
	})
	if !isCreated {
		return
	}

	s.Run("Объявление доступно по ID и возвращает переданные при создании параметры", func() {
		base.Step(fmt.Sprintf("Запрашиваем объявление по id=%s", created.ID))
		got := advertisements.GetAdvertisementByID(s.T(), token, created.ID)

		base.Check("Параметры объявления совпадают с переданными при создании")
		assertions.AdvertisementMatchesRequest(s.T(), got, request, s.loginData.User.Id)
	})

	s.Run("Объявление находится поиском по полному названию", func() {
		base.Step(fmt.Sprintf("Ищем объявления по названию %s", request.Title))
		list := advertisements.GetAdvertisementsList(s.T(), token, request.Title)

		base.Check("Созданное объявление есть в результатах поиска")
		assertions.ContainsAdvertisementByTitle(s.T(), list.Items, request.Title)
	})

	s.Run("Фотография загружена и доступна по ручке фотографий", func() {
		base.Step(fmt.Sprintf("Запрашиваем фотографии объявления id=%s", created.ID))
		photos := advertisements.GetAdvertisementPhotos(s.T(), token, created.ID)

		base.Check("Количество фотографий совпадает с загруженным и у каждой есть url")
		require.Lenf(s.T(), photos, len(request.PhotoPaths), "Количество фотографий не совпадает с загруженным, объявление id=%s", created.ID)
		for _, photo := range photos {
			require.NotEmptyf(s.T(), photo.URL, "У фотографии id=%s пустой url", photo.ID)
		}
	})
}

// Negative: некорректные параметры при создании объявления возвращают 400 и сообщение об ошибке
func (s *TestSuite) TestCreateAdvertisementInvalidParams() {
	cases := []struct {
		name        string
		modify      func(r *advertisementModels.CreateAdvertisementRequest)
		wantMessage string
	}{
		{"пустой title", func(r *advertisementModels.CreateAdvertisementRequest) { r.Title = "" }, "title is required"},
		{"слишком длинный title", func(r *advertisementModels.CreateAdvertisementRequest) {
			r.Title = strings.Repeat("a", maxTitleLength+1)
		}, fmt.Sprintf("title must not exceed %d characters provided %d", serverTitleLimit, maxTitleLength+1)},
		{"отрицательная цена", func(r *advertisementModels.CreateAdvertisementRequest) { r.Price = -100 },
			"price must be between 0 and 1000000 inclusive"},
		{"пустое описание", func(r *advertisementModels.CreateAdvertisementRequest) { r.Description = "" }, "description is required"},
		{"без фотографий", func(r *advertisementModels.CreateAdvertisementRequest) { r.PhotoPaths = nil }, "at least 1 photo is required"},
		{"отрицательное количество", func(r *advertisementModels.CreateAdvertisementRequest) { r.Quantity = negativeQuantity },
			"quantity must be between 1 and 100 inclusive"},
		{"строка вместо цены", func(r *advertisementModels.CreateAdvertisementRequest) { r.PriceRaw = "abc" },
			"cannot convert string value abc to float"},
		{"строка вместо количества", func(r *advertisementModels.CreateAdvertisementRequest) { r.QuantityRaw = "abc" },
			"cannot convert string value to int"},
		{"неподдерживаемый формат файла фотографии", func(r *advertisementModels.CreateAdvertisementRequest) {
			r.PhotoPaths = []string{invalidPhotoPath}
		}, `photo 1: file type not allowed: extension ".txt"`},
		{"больше трёх фотографий", func(r *advertisementModels.CreateAdvertisementRequest) {
			r.PhotoPaths = validPhotos(maxPhotosCount + 1)
		}, fmt.Sprintf("maximum %d photos allowed", maxPhotosCount)},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			request := advertisements.NewAdvertisementRequest(validPhotos(1)...)
			tc.modify(&request)

			base.Step(fmt.Sprintf("Пытаемся создать объявление: %s", tc.name))
			errResponse := advertisements.CreateAdvertisementExpectError(s.T(), s.loginData.Token, request, http.StatusBadRequest)

			base.Check("Вернулся 400, error и message совпадают с ожидаемыми")
			assertions.AssertErrorResponseEquals(s.T(), errResponse, assertions.InvalidRequestErrorCode, tc.wantMessage)
		})
	}
}

// Negative: запрос без авторизации или с неверным токеном возвращает 401 и сообщение об ошибке
func (s *TestSuite) TestCreateAdvertisementUnauthorized() {
	cases := []struct {
		name        string
		token       string
		wantMessage string
	}{
		{"без токена", "", assertions.MissingTokenMessage},
		{"неверный токен", incorrectToken, assertions.InvalidTokenMessage},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			time.Sleep(300 * time.Millisecond)
			base.Step(fmt.Sprintf("Пытаемся создать объявление: %s", tc.name))
			errResponse := advertisements.CreateAdvertisementExpectError(s.T(), tc.token,
				advertisements.NewAdvertisementRequest(validPhotos(1)...), http.StatusUnauthorized)

			base.Check("Вернулся 401 и корректное сообщение об ошибке")
			assertions.AssertUnauthorizedResponse(s.T(), errResponse, tc.wantMessage)
		})
	}
}
