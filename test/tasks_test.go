package test

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sebrupomo/sebrupomo-backend/handler"
	"github.com/sebrupomo/sebrupomo-backend/jwt"
	"github.com/stretchr/testify/assert"
)

func TestCreateTasks(t *testing.T) {
	c, rec := getTestWithAutoAuth(http.MethodPost, "?title=TestTask", "")

	err := jwt.JwtMiddlware(func(context echo.Context) error {
		return handler.CreateTask(c)
	})(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, len(getUser().Tasks), 1)
	}

	taskId := getUser().Tasks[0].ID
	c, rec = getTestWithAutoAuth(http.MethodPost, "api/me/task/:path?title=TestTask", "")
	c.SetParamNames("path")
	c.SetParamValues(taskId)

	err = jwt.JwtMiddlware(func(context echo.Context) error {
		return handler.CreateTask(c)
	})(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, len(getUser().Tasks[0].SubTasks), 1)
	}
}

func TestCreateTaskMissingTitle(t *testing.T) {
	c, rec := getTestWithAutoAuth(http.MethodPost, "", "")

	err := jwt.JwtMiddlware(func(context echo.Context) error {
		return handler.CreateTask(c)
	})(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestRemoveTasks(t *testing.T) {
	if len(getUser().Tasks) == 0 {
		TestCreateTasks(t)
	}

	c, rec := getTestWithAutoAuth(http.MethodDelete, "api/me/task/:path", "")
	c.SetParamNames("path")
	c.SetParamValues(getUser().Tasks[0].ID + "/" + getUser().Tasks[0].SubTasks[0].ID)

	err := jwt.JwtMiddlware(func(context echo.Context) error {
		return handler.RemoveTask(c)
	})(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, len(getUser().Tasks), 1)
		assert.Equal(t, len(getUser().Tasks[0].SubTasks), 0)
	}

	c, rec = getTestWithAutoAuth(http.MethodDelete, "api/me/task/:path", "")
	c.SetParamNames("path")
	c.SetParamValues(getUser().Tasks[0].ID)

	err = jwt.JwtMiddlware(func(context echo.Context) error {
		return handler.RemoveTask(c)
	})(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, len(getUser().Tasks), 0)
	}
}
