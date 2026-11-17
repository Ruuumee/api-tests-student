package advertisementSearch

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

const partLength = 5

// TestSuite - набор тестов для функционала поиска объявлений по названию.
type TestSuite struct {
	suite.Suite
	loginData       authModels.LoginOkResponse
	fixture         advertisements.Fixture
	existingAdTitle string
	searchBeginning string
	searchMiddle    string
	searchEnd       string
}

// точка входа для запуска всех тестов поиска объявлений в этом файле.
func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

// Инициализирует окружение, авторизует пользователя и создает эталонное объявление со случайным названием.
func (s *TestSuite) SetupSuite() {
	time.Sleep(500 * time.Millisecond)
	base.SetupSuite()

	base.Precondition("Авторизация пользователя с кредами из переменных окружения")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))

	base.Precondition("Создаем объявление со случайным названием из трёх частей")
	s.searchBeginning = utils.RandomString(partLength)
	s.searchMiddle = utils.RandomString(partLength)
	s.searchEnd = utils.RandomString(partLength)

	request := advertisements.NewAdvertisementRequest()
	request.Title = s.searchBeginning + s.searchMiddle + s.searchEnd
	s.existingAdTitle = s.fixture.CreateAdvertisement(s.T(), s.loginData.Token, request).Title
	require.NotEmpty(s.T(), s.existingAdTitle, "Название созданного объявления пустое")
}

// Помогает удалить все созданные объявления
func (s *TestSuite) TearDownSuite() {
	s.fixture.Cleanup(s.T())
	base.TearDownSuite()
}

// Positive: поиск по подстроке из середины названия находит наше объявление
func (s *TestSuite) TestSearchByTitle() {
	base.Step(fmt.Sprintf("Ищем объявления по подстроке %s", s.searchMiddle))
	list := advertisements.GetAdvertisementsList(s.T(), s.loginData.Token, s.searchMiddle)

	s.Run("В ответе есть хотя бы одно объявление", func() {
		require.NotEmpty(s.T(), list.Items)
	})

	s.Run("Среди найденных есть наше объявление", func() {
		assertions.ContainsAdvertisementByTitle(s.T(), list.Items, s.existingAdTitle)
	})
}

// Positive: поиск находит объявление по подстроке в начале, середине, конце названия и в разных регистрах
func (s *TestSuite) TestSearchByTitleSubstring() {
	cases := []struct {
		name   string
		search string
	}{
		{"начало", s.searchBeginning},
		{"середина", s.searchMiddle},
		{"конец", s.searchEnd},
		{"нижний регистр", strings.ToLower(s.searchMiddle)},
		{"верхний регистр", strings.ToUpper(s.searchMiddle)},
		{"смешанный регистр", utils.MixCase(s.searchMiddle)},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			time.Sleep(200 * time.Millisecond)
			base.Step(fmt.Sprintf("Ищем объявления по подстроке %s (%s)", tc.search, tc.name))
			list := advertisements.GetAdvertisementsList(s.T(), s.loginData.Token, tc.search)

			base.Check("Наше объявление есть в результатах поиска")
			assertions.ContainsAdvertisementByTitle(s.T(), list.Items, s.existingAdTitle)
		})
	}
}
