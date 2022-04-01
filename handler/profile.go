package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sebrupomo/sebrupomo-backend/db"
	"github.com/sebrupomo/sebrupomo-backend/jwt"
	"github.com/sebrupomo/sebrupomo-backend/model"
)

type userSignUpRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Hash     string `json:"hash" validate:"required"`
}

func SignUp(c echo.Context) error {
	u := new(userSignUpRequest)
	if err := (&Context{c}).BindValidate(u); err != nil {
		if err.Error() == "Key: 'userSignUpRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag" {
			return c.JSON(http.StatusBadRequest, model.InputError("Ungültige Email Adresse.", "Email"))
		}
		return c.JSON(http.StatusBadRequest, model.ProcessError())
	}

	if db.ExistsUserByEmail(u.Email) {
		return c.JSON(http.StatusBadRequest, model.InputError("Diese Email-Adresse ist bereits vergeben.", "Email"))
	}

	if db.ExistsUserByUsername(u.Username) {
		return c.JSON(http.StatusBadRequest, model.InputError("Diese Benutzername ist bereits vergeben.", "Benutzername"))
	}

	result, err := db.CreateUser(&model.User{
		Username: u.Username,
		Email:    u.Email,
		Hash:     u.Hash,
		Tasks:    []model.Task{},
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.InternalError())
	}

	return c.JSON(http.StatusCreated, result)
}

type userLoginRequest struct {
	Username string `query:"username" validate:"required"`
	Hash     string `query:"hash" validate:"required"`
}

func Login(c echo.Context) error {
	u := new(userLoginRequest)
	if err := (&Context{c}).BindValidateParams(u); err != nil {
		return c.JSON(http.StatusBadRequest, model.ProcessError())
	}

	user, err := db.FindUserByUsername(u.Username)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.InputError("Benutzer wurde nicht gefunden.", "Benutzername"))
	}
	if user.Hash != u.Hash {
		return c.JSON(http.StatusBadRequest, model.InputError("Falsches Passwort.", "Passwort"))
	}

	token, err := jwt.GenerateToken(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.InternalError())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

func WhoAmI(c echo.Context) error {
	user, err := jwt.GetJwtUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, model.AccessError())
	}
	return c.JSON(http.StatusOK, user)
}
