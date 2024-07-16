package controllers

import (
	"database/sql"
	"errors"
	"github.com/drTragger/MykroTask/middleware"
	"github.com/drTragger/MykroTask/models"
	"github.com/drTragger/MykroTask/services"
	"github.com/drTragger/MykroTask/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

type TaskController struct {
	taskService services.TaskService
}

func NewTaskController(taskService services.TaskService) *TaskController {
	return &TaskController{taskService: taskService}
}

func (tc *TaskController) CreateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectIdStr, ok := vars["projectId"]
	if !ok {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Missing projectId parameter.",
		})
		return
	}

	projectId, err := uuid.Parse(projectIdStr)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Invalid projectId parameter.",
		})
		return
	}

	var task *models.Task
	errorResponse := utils.UnmarshalRequest(r, &task)
	if errorResponse != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, errorResponse)
		return
	}

	// Validate the request data
	err = utils.ValidateStruct(task)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusUnprocessableEntity, &utils.ErrorResponse{
			Status:  false,
			Message: "Validation failed.",
			Errors:  err.Error(),
		})
		return
	}

	userId := uuid.MustParse(r.Context().Value(middleware.UserIDKey).(string))
	task.ProjectID = projectId
	task.CreatedBy = userId

	task, forbidden, err := tc.taskService.CreateTask(task)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, &utils.ErrorResponse{
			Status:  false,
			Message: "Failed to create task.",
			Errors:  err.Error(),
		})
		return
	}
	if forbidden {
		utils.WriteJSONResponse(w, http.StatusForbidden, &utils.ErrorResponse{
			Status:  false,
			Message: "You are not a member of this project.",
		})
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, &utils.SuccessResponse{
		Status:  true,
		Message: "Task created successfully.",
		Data:    task,
	})
}

func (tc *TaskController) GetTasksForUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectIdStr, ok := vars["projectId"]
	if !ok {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Missing projectId parameter.",
		})
		return
	}

	projectId, err := uuid.Parse(projectIdStr)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Invalid projectId parameter.",
			Errors:  err.Error(),
		})
		return
	}

	memberIdStr, ok := vars["memberId"]
	if !ok {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Missing memberId parameter.",
		})
		return
	}

	memberId, err := uuid.Parse(memberIdStr)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Invalid memberId parameter.",
			Errors:  err.Error(),
		})
		return
	}

	userId := uuid.MustParse(r.Context().Value(middleware.UserIDKey).(string))

	tasks, forbidden, err := tc.taskService.GetTasksForUser(projectId, memberId, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSONResponse(w, http.StatusNotFound, &utils.ErrorResponse{
				Status:  false,
				Message: "No tasks found.",
			})
			return
		}
		utils.WriteJSONResponse(w, http.StatusInternalServerError, &utils.ErrorResponse{
			Status:  false,
			Message: "Failed to get tasks for user.",
			Errors:  err.Error(),
		})
		return
	}
	if forbidden {
		utils.WriteJSONResponse(w, http.StatusForbidden, &utils.ErrorResponse{
			Status:  false,
			Message: "You are not a member of this project.",
		})
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, &utils.SuccessResponse{
		Status:  true,
		Message: "Tasks retrieved successfully.",
		Data:    tasks,
	})
}

func (tc *TaskController) GetTaskById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectIdStr, ok := vars["projectId"]
	if !ok {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Missing projectId parameter.",
		})
		return
	}
	projectId, err := uuid.Parse(projectIdStr)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Invalid projectId parameter.",
			Errors:  err.Error(),
		})
		return
	}
	taskIdStr, ok := vars["taskId"]
	if !ok {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Missing taskId parameter.",
		})
		return
	}
	taskId, err := uuid.Parse(taskIdStr)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Invalid taskId parameter.",
			Errors:  err.Error(),
		})
		return
	}
	userId := uuid.MustParse(r.Context().Value(middleware.UserIDKey).(string))

	task, forbidden, err := tc.taskService.GetTaskById(projectId, taskId, userId)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, &utils.ErrorResponse{
			Status:  false,
			Message: "Failed to get task.",
			Errors:  err.Error(),
		})
		return
	}
	if forbidden {
		utils.WriteJSONResponse(w, http.StatusForbidden, &utils.ErrorResponse{
			Status:  false,
			Message: "You are not a member of this project.",
		})
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, &utils.SuccessResponse{
		Status:  true,
		Message: "Task retrieved successfully.",
		Data:    task,
	})
}

func (tc *TaskController) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectIdStr, ok := vars["projectId"]
	if !ok {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Missing projectId parameter.",
		})
		return
	}
	projectId, err := uuid.Parse(projectIdStr)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Invalid projectId parameter.",
			Errors:  err.Error(),
		})
	}
	taskIdStr, ok := vars["taskId"]
	if !ok {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Missing taskId parameter.",
		})
		return
	}
	taskId, err := uuid.Parse(taskIdStr)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Invalid taskId parameter.",
			Errors:  err.Error(),
		})
		return
	}
	userId := uuid.MustParse(r.Context().Value(middleware.UserIDKey).(string))

	forbidden, err := tc.taskService.DeleteTask(projectId, userId, taskId)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, &utils.ErrorResponse{
			Status:  false,
			Message: "Failed to delete task.",
			Errors:  err.Error(),
		})
		return
	}
	if forbidden {
		utils.WriteJSONResponse(w, http.StatusForbidden, &utils.ErrorResponse{
			Status:  false,
			Message: "You are not allowed to edit this project.",
		})
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, &utils.SuccessResponse{
		Status:  true,
		Message: "Task deleted successfully.",
	})
}

func (tc *TaskController) GetTasksForProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectIdStr, ok := vars["projectId"]
	if !ok {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Missing projectId parameter.",
		})
		return
	}
	projectId, err := uuid.Parse(projectIdStr)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Invalid projectId parameter.",
			Errors:  err.Error(),
		})
		return
	}
	userId := uuid.MustParse(r.Context().Value(middleware.UserIDKey).(string))

	tasks, forbidden, err := tc.taskService.GetTasksForProject(projectId, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSONResponse(w, http.StatusNotFound, &utils.ErrorResponse{
				Status:  false,
				Message: "No tasks found for project.",
			})
			return
		}
		utils.WriteJSONResponse(w, http.StatusInternalServerError, &utils.ErrorResponse{
			Status:  false,
			Message: "Failed to get tasks for project.",
			Errors:  err.Error(),
		})
		return
	}
	if forbidden {
		utils.WriteJSONResponse(w, http.StatusForbidden, &utils.ErrorResponse{
			Status:  false,
			Message: "You are not a member of this project.",
		})
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, &utils.SuccessResponse{
		Status:  true,
		Message: "Tasks retrieved successfully.",
		Data:    tasks,
	})
}

func (tc *TaskController) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectIdStr, ok := vars["projectId"]
	if !ok {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Missing projectId parameter.",
		})
		return
	}
	projectId, err := uuid.Parse(projectIdStr)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Invalid projectId parameter.",
			Errors:  err.Error(),
		})
		return
	}
	taskIdStr, ok := vars["taskId"]
	if !ok {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Missing taskId parameter.",
		})
		return
	}
	taskId, err := uuid.Parse(taskIdStr)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Invalid taskId parameter.",
			Errors:  err.Error(),
		})
		return
	}
	userId := uuid.MustParse(r.Context().Value(middleware.UserIDKey).(string))

	var task *models.Task
	errorResponse := utils.UnmarshalRequest(r, &task)
	if errorResponse != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, errorResponse)
		return
	}

	// Validate the request data
	err = utils.ValidateStruct(task)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusUnprocessableEntity, &utils.ErrorResponse{
			Status:  false,
			Message: "Validation failed.",
			Errors:  err.Error(),
		})
		return
	}

	task, forbidden, err := tc.taskService.UpdateTask(projectId, taskId, userId, task)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSONResponse(w, http.StatusNotFound, &utils.ErrorResponse{
				Status:  false,
				Message: "No task found for project.",
			})
			return
		}
		utils.WriteJSONResponse(w, http.StatusInternalServerError, &utils.ErrorResponse{
			Status:  false,
			Message: "Failed to update task.",
			Errors:  err.Error(),
		})
		return
	}
	if forbidden {
		utils.WriteJSONResponse(w, http.StatusForbidden, &utils.ErrorResponse{
			Status:  false,
			Message: "You are not allowed to edit this project.",
		})
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, &utils.SuccessResponse{
		Status:  true,
		Message: "Task updated successfully.",
		Data:    task,
	})
}
