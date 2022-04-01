package test

import (
	"encoding/json"
	"net/http/httptest"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/sebrupomo/sebrupomo-backend/db"
	"github.com/sebrupomo/sebrupomo-backend/handler"
	"github.com/sebrupomo/sebrupomo-backend/jwt"
	"github.com/sebrupomo/sebrupomo-backend/model"
)

func setupDatabase() {
	if db.Database == nil {
		db.GetConnection(true)
	}
}

func setupEnvironment() *echo.Echo {
	setupDatabase()

	e := echo.New()
	e.Validator = &handler.CustomValidator{Validator: validator.New()}
	return e
}

func getTestWithAuth(method string, query string, body string, auth string) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, strings.ReplaceAll("/"+query, " ", "%20"), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+auth)
	rec := httptest.NewRecorder()
	return setupEnvironment().NewContext(req, rec), rec
}

func getTestWithAutoAuth(method string, query string, body string) (echo.Context, *httptest.ResponseRecorder) {
	token, _ := getToken()
	return getTestWithAuth(method, query, body, token)
}

func getTest(method string, query string, body string) (echo.Context, *httptest.ResponseRecorder) {
	return getTestWithAuth(method, query, body, "")
}

func responseMap(b []byte) map[string]interface{} {
	var m map[string]interface{}
	json.Unmarshal(b, &m)
	return m
}

func createUser() *model.User {
	setupDatabase()

	user := &model.User{
		Username: "John Doe",
		Email:    "test@web.de",
		Hash:     "test123",
		Tasks:    []model.Task{},
	}

	db.CreateUser(user)

	return user
}

func getUser() *model.User {
	setupDatabase()

	user, err := db.FindUserByUsername("John Doe")
	if err != nil {
		return createUser()
	}
	return user
}

func getToken() (string, error) {
	return jwt.GenerateToken(getUser())
}
