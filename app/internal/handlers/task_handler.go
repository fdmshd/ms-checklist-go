package handlers

import (
	"checklist/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type TaskHandler struct {
	TaskModel models.TaskModel
}

func (h *TaskHandler) Get(c echo.Context) error {
	return nil
}

func (h *TaskHandler) GetAll(c echo.Context) error {
	requestedID := c.Param("user_id")
	requester := userIDFromToken(c)

	if !(requestedID == requester || isAdminFromToken(c)) {
		return &echo.HTTPError{Code: http.StatusForbidden, Message: "forbidden"}
	}
	id, err := strconv.Atoi(requestedID)
	if err != nil {
		panic("something went wrong with user_id")
	}
	tasks, err := h.TaskModel.GetByUser(id)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	return c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) Create(c echo.Context) error {
	requesterID := userIDFromToken(c)
	id, err := strconv.Atoi(requesterID)
	if err != nil {
		panic("something went wrong with user_id")
	}
	task := new(models.Task)
	if err := c.Bind(task); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}
	task.UserID = id
	task.CreatedAt = time.Now().Format("2006-01-02T15:04:05")

	err = c.Validate(task)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	taskID, err := h.TaskModel.Add(task)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}
	task.Id = taskID
	return c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) Update(c echo.Context) error {
	currUserID := userIDFromToken(c)
	TaskId := c.Param("id")
	oldTask, err := h.getTask(currUserID, TaskId)
	if err != nil {
		return err
	}
	task := new(models.Task)
	if err = c.Bind(task); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}
	task.Id = oldTask.Id
	err = c.Validate(task)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	err = h.TaskModel.Update(task)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}
	return c.JSON(http.StatusOK, "updated")
}

func (h *TaskHandler) Delete(c echo.Context) error {
	currUserID := userIDFromToken(c)
	TaskId := c.Param("id")
	task, err := h.getTask(currUserID, TaskId)
	if err != nil {
		return err
	}
	err = h.TaskModel.Delete(task.Id)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}
	return c.JSON(http.StatusOK, "deleted")
}

func (h *TaskHandler) Complete(c echo.Context) error {
	currUserID := userIDFromToken(c)
	TaskId := c.Param("id")
	task, err := h.getTask(currUserID, TaskId)
	if err != nil {
		return err
	}
	err = h.TaskModel.UpdateFields(map[string]interface{}{"is_done": !task.IsDone}, task.Id)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}
	if !task.IsDone {
		return c.JSON(http.StatusOK, "completed")
	}
	return c.JSON(http.StatusOK, "uncompleted")
}

func (h *TaskHandler) getTask(currUserID, TaskId string) (*models.Task, error) {

	id, err := strconv.Atoi(TaskId)
	if err != nil {
		panic("something went wrong with id")
	}
	task, err := h.TaskModel.Get(id)
	if err != nil {
		return nil, &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}
	userID, err := strconv.Atoi(currUserID)
	if err != nil {
		panic("something went wrong with user_id")
	}
	if userID != task.UserID {
		return nil, &echo.HTTPError{Code: http.StatusForbidden, Message: "forbidden"}
	}
	task.Id = id
	return task, nil
}

func userIDFromToken(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["id"].(string)
}

func isAdminFromToken(c echo.Context) bool {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["is_admin"].(bool)
}
