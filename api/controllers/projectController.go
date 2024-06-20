package controllers

import (
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
			Message: "Missing id parameter",
		})
		return
	}

	projectId, err := uuid.Parse(idStr)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Invalid id parameter",
		})
		return
	}

	userId := uuid.MustParse(r.Context().Value(middleware.UserIDKey).(string))

	project, err := pc.projectService.GetProjectById(projectId)
	if err != nil {
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
