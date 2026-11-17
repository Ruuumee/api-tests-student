package favorites

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/advertisements"
	advertisementModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

const incorrectToken = "incorrect_token"

// TestSuite - набор тестов для функционала "Избранное"
// Содержит двух пользователей: покупателя и продавца, а также фикстуру для управления объявлениями
type TestSuite struct {
	suite.Suite
	buyer   authModels.LoginOkResponse
	seller  authModels.LoginOkResponse
	fixture advertisements.Fixture
}

// TestSuiteRun - точка входа для запуска всех тестов в этом файле
func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

// Инициализирует тестовое окружение и создает двух пользователей
func (s *TestSuite) SetupSuite() {
	time.Sleep(1 * time.Second)
	base.SetupSuite()

	base.Precondition("Авторизация покупателя с кредами из переменных окружения")
	s.buyer = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))

	base.Precondition("Регистрация нового пользователя-продавца")
	s.seller = auth.RegisterRandomUser(s.T())
}

// Очищает все объявления, созданные в рамках текущего теста
func (s *TestSuite) TearDownTest() {
	s.fixture.Cleanup(s.T())
}

// Выполняет финальную очистку тестового окружения
func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// createSellerAdvertisement создает объявление от имени продавца, после теста оно удалится в TearDownTest
func (s *TestSuite) createSellerAdvertisement() advertisementModels.Advertisement {
	base.Precondition("Продавец создает объявление")
	return s.fixture.CreateAdvertisement(s.T(), s.seller.Token, advertisements.NewAdvertisementRequest())
}

// Positive: добавленное в избранное объявление появляется в списке, удалённое - пропадает
func (s *TestSuite) TestFavoritesFullCycle() {
	ad := s.createSellerAdvertisement()

	s.Run("Добавляем объявление в избранное", func() {
		time.Sleep(300 * time.Millisecond)
		base.Step(fmt.Sprintf("Покупатель добавляет в избранное объявление id=%s", ad.ID))
		favorite := favorites.AddToFavorites(s.T(), s.buyer.Token, ad.ID)

		base.Check("В ответе вернулось именно это объявление")
		require.Equal(s.T(), ad.ID, favorite.Advertisement.ID)
	})

	s.Run("Объявление появилось в списке избранного", func() {
		time.Sleep(300 * time.Millisecond)
		base.Step("Запрашиваем список избранного")
		list := favorites.GetFavorites(s.T(), s.buyer.Token)

		base.Check("Объявление есть в списке")
		assertions.FavoritesContainAdvertisement(s.T(), list.Items, ad.ID)
	})

	s.Run("Удаляем объявление из избранного", func() {
		time.Sleep(300 * time.Millisecond)
		base.Step(fmt.Sprintf("Покупатель удаляет из избранного объявление id=%s", ad.ID))
		favorites.RemoveFromFavorites(s.T(), s.buyer.Token, ad.ID)
	})

	s.Run("Объявление пропало из списка избранного", func() {
		time.Sleep(300 * time.Millisecond)
		base.Step("Запрашиваем список избранного")
		list := favorites.GetFavorites(s.T(), s.buyer.Token)

		base.Check("Объявления нет в списке")
		assertions.FavoritesNotContainAdvertisement(s.T(), list.Items, ad.ID)
	})
}

// Positive: если удалить само объявление, оно исчезает из избранного пользователя
func (s *TestSuite) TestFavoriteDisappearsAfterAdvertisementDeleted() {
	ad := s.createSellerAdvertisement()

	base.Precondition("Покупатель добавляет объявление в избранное")
	favorites.AddToFavorites(s.T(), s.buyer.Token, ad.ID)
	assertions.FavoritesContainAdvertisement(s.T(), favorites.GetFavorites(s.T(), s.buyer.Token).Items, ad.ID)

	s.Run("Продавец удаляет объявление", func() {
		base.Step(fmt.Sprintf("Удаляем объявление id=%s", ad.ID))
		advertisements.DeleteAdvertisement(s.T(), s.seller.Token, ad.ID)
	})

	s.Run("Объявление пропало из избранного покупателя", func() {
		base.Step("Запрашиваем список избранного")
		list := favorites.GetFavorites(s.T(), s.buyer.Token)

		base.Check("Удалённого объявления нет в списке")
		assertions.FavoritesNotContainAdvertisement(s.T(), list.Items, ad.ID)
	})
}

// Negative: добавление в избранное несуществующего объявления возвращает 404
func (s *TestSuite) TestAddNotExistingAdvertisement() {
	notExistingID := utils.RandomUUID()

	base.Step(fmt.Sprintf("Добавляем в избранное несуществующее объявление id=%s", notExistingID))
	errResponse := favorites.AddToFavoritesExpectError(s.T(), s.buyer.Token, notExistingID, http.StatusNotFound)

	base.Check("Вернулся 404 и заполнены поля error и message")
	assertions.AssertErrorResponse(s.T(), errResponse)
}

// Negative: удаление из избранного несуществующего объявления возвращает 404
func (s *TestSuite) TestRemoveNotExistingAdvertisement() {
	notExistingID := utils.RandomUUID()

	base.Step(fmt.Sprintf("Удаляем из избранного несуществующее объявление id=%s", notExistingID))
	errResponse := favorites.RemoveFromFavoritesExpectError(s.T(), s.buyer.Token, notExistingID, http.StatusNotFound)

	base.Check("Вернулся 404 и заполнены поля error и message")
	assertions.AssertErrorResponse(s.T(), errResponse)
}

// Negative: ручки избранного без токена или с неверным токеном возвращают 401
func (s *TestSuite) TestFavoritesUnauthorized() {
	adID := utils.RandomUUID()

	operations := []struct {
		name string
		call func(t *testing.T, token string) response.ErrorResponse
	}{
		{"POST /favorites", func(t *testing.T, token string) response.ErrorResponse {
			return favorites.AddToFavoritesExpectError(t, token, adID, http.StatusUnauthorized)
		}},
		{"GET /favorites", func(t *testing.T, token string) response.ErrorResponse {
			return favorites.GetFavoritesExpectError(t, token, http.StatusUnauthorized)
		}},
		{"DELETE /favorites/{ad_id}", func(t *testing.T, token string) response.ErrorResponse {
			return favorites.RemoveFromFavoritesExpectError(t, token, adID, http.StatusUnauthorized)
		}},
	}
	tokens := []struct {
		name  string
		value string
	}{
		{"без токена", ""},
		{"неверный токен", incorrectToken},
	}

	for _, operation := range operations {
		for _, token := range tokens {
			s.Run(operation.name+": "+token.name, func() {
				base.Step(fmt.Sprintf("Вызываем %s (%s)", operation.name, token.name))
				errResponse := operation.call(s.T(), token.value)

				base.Check("Вернулся 401 и корректное сообщение об ошибке")
				assertions.AssertUnauthorizedResponse(s.T(), errResponse)
			})
		}
	}
}

// Negative: повторное добавление того же объявления возвращает 409
func (s *TestSuite) TestAddAlreadyFavoriteAdvertisement() {
	ad := s.createSellerAdvertisement()

	base.Precondition("Покупатель добавляет объявление в избранное")
	favorites.AddToFavorites(s.T(), s.buyer.Token, ad.ID)

	base.Step(fmt.Sprintf("Повторно добавляем объявление id=%s", ad.ID))
	errResponse := favorites.AddToFavoritesExpectError(s.T(), s.buyer.Token, ad.ID, http.StatusConflict)

	base.Check("Вернулся 409 и заполнены поля error и message")
	assertions.AssertErrorResponse(s.T(), errResponse)
}

// Negative: нельзя добавить в избранное собственное объявление (403)
func (s *TestSuite) TestAddOwnAdvertisement() {
	base.Precondition("Покупатель создает собственное объявление")
	ownAd := s.fixture.CreateAdvertisement(s.T(), s.buyer.Token, advertisements.NewAdvertisementRequest())

	base.Step(fmt.Sprintf("Добавляем в избранное собственное объявление id=%s", ownAd.ID))
	errResponse := favorites.AddToFavoritesExpectError(s.T(), s.buyer.Token, ownAd.ID, http.StatusForbidden)

	base.Check("Вернулся 403 и заполнены поля error и message")
	assertions.AssertErrorResponse(s.T(), errResponse)
}

// Negative: запрос на добавление в избранное с пустым advertisement_id возвращает 400
func (s *TestSuite) TestAddWithEmptyAdvertisementID() {
	base.Step("Добавляем в избранное объявление с пустым advertisement_id")
	errResponse := favorites.AddToFavoritesExpectError(s.T(), s.buyer.Token, "", http.StatusBadRequest)

	base.Check("Вернулся 400 и заполнены поля error и message")
	assertions.AssertErrorResponse(s.T(), errResponse)
}
