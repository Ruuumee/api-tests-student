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
	incorrectToken = "incorrect_token"
)

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
	request := advertisements.NewAdvertisementRequest()

	var created advertisementModels.Advertisement
	isCreated := s.Run("Создаем объявление со всеми параметрами и фотографией", func() {
		base.Step(fmt.Sprintf("Создаем объявление с названием %s", request.Title))
		created = s.fixture.CreateAdvertisement(s.T(), token, request)

		base.Check("Объявление создано и в ответе есть id")
		require.NotEmpty(s.T(), created.ID)
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
		require.Len(s.T(), photos, request.PhotosCount, "Количество фотографий не совпадает с загруженным")
		for _, photo := range photos {
			require.NotEmpty(s.T(), photo.URL, "У фотографии пустой url")
		}
	})
}

// Negative: некорректные параметры при создании объявления возвращают 400 и сообщение об ошибке
func (s *TestSuite) TestCreateAdvertisementInvalidParams() {
	cases := []struct {
		name   string
		modify func(r *advertisementModels.CreateAdvertisementRequest)
	}{
		{"пустой title", func(r *advertisementModels.CreateAdvertisementRequest) { r.Title = "" }},
		{"слишком длинный title", func(r *advertisementModels.CreateAdvertisementRequest) {
			r.Title = strings.Repeat("a", maxTitleLength+1)
		}},
		{"отрицательная цена", func(r *advertisementModels.CreateAdvertisementRequest) { r.Price = -100 }},
		{"пустое описание", func(r *advertisementModels.CreateAdvertisementRequest) { r.Description = "" }},
		{"без фотографий", func(r *advertisementModels.CreateAdvertisementRequest) { r.PhotosCount = 0 }},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			request := advertisements.NewAdvertisementRequest()
			tc.modify(&request)

			base.Step(fmt.Sprintf("Пытаемся создать объявление: %s", tc.name))
			errResponse := advertisements.CreateAdvertisementExpectError(s.T(), s.loginData.Token, request, http.StatusBadRequest)

			base.Check("Вернулся 400 и заполнены поля error и message")
			assertions.AssertErrorResponse(s.T(), errResponse)
		})
	}
}

// Negative: запрос без авторизации или с неверным токеном возвращает 401 и сообщение об ошибке
func (s *TestSuite) TestCreateAdvertisementUnauthorized() {
	cases := []struct {
		name  string
		token string
	}{
		{"без токена", ""},
		{"неверный токен", incorrectToken},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			time.Sleep(300 * time.Millisecond)
			base.Step(fmt.Sprintf("Пытаемся создать объявление: %s", tc.name))
			errResponse := advertisements.CreateAdvertisementExpectError(s.T(), tc.token,
				advertisements.NewAdvertisementRequest(), http.StatusUnauthorized)

			base.Check("Вернулся 401 и корректное сообщение об ошибке")
			assertions.AssertUnauthorizedResponse(s.T(), errResponse)
		})
	}
}
