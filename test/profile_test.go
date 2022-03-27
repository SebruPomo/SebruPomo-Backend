package test

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sebrupomo/sebrupomo-backend/db"
	"github.com/sebrupomo/sebrupomo-backend/handler"
	"github.com/sebrupomo/sebrupomo-backend/jwt"
	"github.com/stretchr/testify/assert"
)

func TestSignUp(t *testing.T) {
	c, rec := getTest(http.MethodPost, "", `{"username":"John Doe","email":"test@gmail.com","hash":"test123"}`)

	if assert.NoError(t, handler.SignUp(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestSignUpUsernameTaken(t *testing.T) {
	createUser()

	c, rec := getTest(http.MethodPost, "", `{"username":"John Doe","email":"test@gmail.com","hash":"test123"}`)
	if assert.NoError(t, handler.SignUp(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestSignUpEmailTaken(t *testing.T) {
	c, rec := getTest(http.MethodPost, "", `{"username":"John Doe","email":"test@gmail.com","hash":"test123"}`)
	if assert.NoError(t, handler.SignUp(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestSignUpInvalidBody(t *testing.T) {
	c, rec := getTest(http.MethodPost, "", `{"username":"John Doe","email":"test@gmail.com"}`)
	if assert.NoError(t, handler.SignUp(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}

	c, rec = getTest(http.MethodPost, "", `{"username":"John Doe","email":"testgmail.com","hash":"test123"}`)
	if assert.NoError(t, handler.SignUp(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}

	c, rec = getTest(http.MethodPost, "", "")
	if assert.NoError(t, handler.SignUp(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestLogin(t *testing.T) {
	createUser()

	c, rec := getTest(http.MethodPost, "?username=John Doe&hash=test123", "")
	if assert.NoError(t, handler.Login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NotEmpty(t, responseMap(rec.Body.Bytes())["token"])
	}
}

func TestLoginInvalidUser(t *testing.T) {
	db.ResetDatabase()

	c, rec := getTest(http.MethodPost, "?username=John Doe&hash=test123", "")
	if assert.NoError(t, handler.Login(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	createUser()

	c, rec := getTest(http.MethodPost, "?username=John Doe&hash=test1234", "")
	if assert.NoError(t, handler.Login(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestLoginInvalidQuery(t *testing.T) {
	c, rec := getTest(http.MethodPost, "?hash=test123", "")
	if assert.NoError(t, handler.Login(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}

	c, rec = getTest(http.MethodPost, "?username=John Doe", "")
	if assert.NoError(t, handler.Login(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestWhoAmI(t *testing.T) {
	token, err := getToken()

	if assert.Empty(t, err) {
		c, rec := getTestWithAuth(http.MethodGet, "api/me/whoAmI", "", token)

		err := jwt.JwtMiddlware(func(context echo.Context) error {
			return handler.WhoAmI(c)
		})(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "John Doe", responseMap(rec.Body.Bytes())["username"])
		}
	}
}

func TestWhoAmIMissingHeader(t *testing.T) {
	c, _ := getTest(http.MethodGet, "api/me/whoAmI", "")

	err := jwt.JwtMiddlware(func(context echo.Context) error {
		return handler.WhoAmI(c)
	})(c)

	assert.Error(t, err)
}

func TestWhoAmIWrongHeader(t *testing.T) {
	c, _ := getTest(http.MethodGet, "api/me/whoAmI", "test123")

	err := jwt.JwtMiddlware(func(context echo.Context) error {
		return handler.WhoAmI(c)
	})(c)

	assert.Error(t, err)
}
