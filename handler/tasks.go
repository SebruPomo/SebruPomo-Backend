package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sebrupomo/sebrupomo-backend/db"
	"github.com/sebrupomo/sebrupomo-backend/jwt"
	"github.com/sebrupomo/sebrupomo-backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createTaskRequest struct {
	Title string `query:"title" validate:"required"`
}

func CreateTask(c echo.Context) error {
	user, err := jwt.GetJwtUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, model.AccessError())
	}

	t := new(createTaskRequest)
	if err := (&Context{c}).BindValidateParams(t); err != nil {
		return c.JSON(http.StatusBadRequest, model.ProcessError())
	}

	task := model.Task{
		ID:          primitive.NewObjectID().Hex(),
		Title:       t.Title,
		Description: "",
		SubTasks:    []model.Task{},
	}

	user.Tasks = addSubTask(task, user.Tasks, c.Param("path"))
	db.UpdateUser(user)

	return c.JSON(http.StatusCreated, task)
}

func addSubTask(task model.Task, subTasks []model.Task, path string) []model.Task {
	if strings.HasPrefix(path, "/") {
		path = strings.Replace(path, "/", "", 1)
	}
	if path == "" {
		subTasks = append(subTasks, task)
	}
	for i, subTask := range subTasks {
		if strings.HasPrefix(path, subTask.ID) {
			subTasks[i].SubTasks = addSubTask(task, subTask.SubTasks, strings.Replace(path, subTask.ID, "", 1))
		}
	}

	return subTasks
}

func RemoveTask(c echo.Context) error {
	user, err := jwt.GetJwtUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, model.AccessError())
	}

	user.Tasks = removeSubTask(user.Tasks, c.Param("path"))
	db.UpdateUser(user)

	return c.NoContent(http.StatusOK)
}

func removeSubTask(subTasks []model.Task, path string) []model.Task {
	if strings.HasPrefix(path, "/") {
		path = strings.Replace(path, "/", "", 1)
	}

	for i, subTask := range subTasks {
		if strings.HasPrefix(path, subTask.ID) {
			path = strings.Replace(path, subTask.ID, "", 1)
			if path == "" {
				subTasks = append(subTasks[:i], subTasks[i+1:]...)
			} else {
				subTasks[i].SubTasks = removeSubTask(subTask.SubTasks, path)
			}
		}
	}

	return subTasks
}
