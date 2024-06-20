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
	"strconv"
)

type ProjectController struct {
	projectService services.ProjectService
	userService    services.UserService
}

func NewProjectController(projectService services.ProjectService, userService services.UserService) *ProjectController {
	return &ProjectController{
		projectService: projectService,
		userService:    userService,
	}
}

func (pc *ProjectController) CreateProject(w http.ResponseWriter, r *http.Request) {
	var project *models.Project
	errorResponse := utils.UnmarshalRequest(r, &project)
	if errorResponse != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, errorResponse)
		return
	}

	// Validate the request data
	err := utils.ValidateStruct(project)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusUnprocessableEntity, &utils.ErrorResponse{
			Status:  false,
			Message: "Validation failed.",
			Errors:  err.Error(),
		})
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	project.OwnerId = uuid.MustParse(userID)

	user, err := pc.userService.GetUserById(project.OwnerId)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, &utils.ErrorResponse{
			Status:  false,
			Message: "Failed to get user.",
			Errors:  err.Error(),
		})
		return
	}
	if user == nil {
		utils.WriteJSONResponse(w, http.StatusUnprocessableEntity, &utils.ErrorResponse{
			Status:  false,
			Message: "User not found.",
		})
		return
	}

	project, err = pc.projectService.CreateProject(project)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, &utils.ErrorResponse{
			Status:  false,
			Message: "Failed to create project.",
			Errors:  err.Error(),
		})
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, &utils.SuccessResponse{
		Status:  true,
		Message: "Project created successfully.",
		Data:    project,
	})
}

func (pc *ProjectController) GetProjectsForUser(w http.ResponseWriter, r *http.Request) {
	userId := uuid.MustParse(r.Context().Value(middleware.UserIDKey).(string))
	var page = 0
	pageParam := r.URL.Query().Get("page")
	if pageParam != "" {
		var err error
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
				Status:  false,
				Message: "Wrong page param.",
				Errors:  err.Error(),
			})
			return
		}
		page--
		if page < 0 {
			page = 0
		}
	}

	projects, err := pc.projectService.GetProjectsForUser(userId, uint(page))
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, &utils.ErrorResponse{
			Status:  false,
			Message: "Failed to get projects.",
			Errors:  err.Error(),
		})
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, &utils.SuccessResponse{
		Status:  true,
		Message: "Projects retrieved successfully.",
		Data:    projects,
	})
}

func (pc *ProjectController) GetProjectById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Missing id parameter.",
		})
		return
	}

	projectId, err := uuid.Parse(idStr)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Invalid id parameter.",
		})
		return
	}

	userId := uuid.MustParse(r.Context().Value(middleware.UserIDKey).(string))

	project, err := pc.projectService.GetProjectById(projectId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSONResponse(w, http.StatusNotFound, &utils.ErrorResponse{
				Status:  false,
				Message: "Project not found.",
			})
			return
		}
		utils.WriteJSONResponse(w, http.StatusInternalServerError, &utils.ErrorResponse{
			Status:  false,
			Message: "Failed to get project.",
			Errors:  err.Error(),
		})
		return
	}

	if project.OwnerId != userId {
		utils.WriteJSONResponse(w, http.StatusForbidden, &utils.ErrorResponse{
			Status:  false,
			Message: "You are not the owner of this project.",
		})
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, &utils.SuccessResponse{
		Status:  true,
		Message: "Project retrieved successfully.",
		Data:    project,
	})
}

func (pc *ProjectController) UpdateProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Missing id parameter.",
		})
		return
	}

	projectId, err := uuid.Parse(idStr)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Invalid id parameter.",
		})
		return
	}

	var project *models.Project
	errorResponse := utils.UnmarshalRequest(r, &project)
	if errorResponse != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, errorResponse)
		return
	}

	err = utils.ValidateStruct(project)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusUnprocessableEntity, &utils.ErrorResponse{
			Status:  false,
			Message: "Validation failed.",
			Errors:  err.Error(),
		})
		return
	}

	project.ID = projectId
	userId := uuid.MustParse(r.Context().Value(middleware.UserIDKey).(string))

	permission, err := pc.projectService.CheckUserPermission(projectId, userId)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, &utils.ErrorResponse{
			Status:  false,
			Message: "Permission check failed.",
			Errors:  err.Error(),
		})
		return
	}
	if !permission {
		utils.WriteJSONResponse(w, http.StatusForbidden, &utils.ErrorResponse{
			Status:  false,
			Message: "You are not the owner of this project.",
		})
		return
	}

	project, err = pc.projectService.UpdateProject(project)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, &utils.ErrorResponse{
			Status:  false,
			Message: "Failed to update project.",
			Errors:  err.Error(),
		})
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, &utils.SuccessResponse{
		Status:  true,
		Message: "Project updated successfully.",
		Data:    project,
	})
}

func (pc *ProjectController) DeleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Missing id parameter.",
		})
		return
	}

	projectId, err := uuid.Parse(idStr)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Invalid id parameter.",
			Errors:  err.Error(),
		})
		return
	}

	userId := uuid.MustParse(r.Context().Value(middleware.UserIDKey).(string))
	permission, err := pc.projectService.CheckUserPermission(projectId, userId)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, &utils.ErrorResponse{
			Status:  false,
			Message: "Permission check failed.",
			Errors:  err.Error(),
		})
		return
	}
	if !permission {
		utils.WriteJSONResponse(w, http.StatusForbidden, &utils.ErrorResponse{
			Status:  false,
			Message: "You are not the owner of this project.",
		})
		return
	}

	err = pc.projectService.DeleteProject(projectId)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, &utils.ErrorResponse{
			Status:  false,
			Message: "Failed to delete project.",
			Errors:  err.Error(),
		})
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, &utils.SuccessResponse{
		Status:  true,
		Message: "Project deleted successfully.",
	})
}
