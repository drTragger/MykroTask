package controllers

import (
	"github.com/drTragger/MykroTask/middleware"
	"github.com/drTragger/MykroTask/models"
	"github.com/drTragger/MykroTask/services"
	"github.com/drTragger/MykroTask/utils"
	"github.com/google/uuid"
	"log"
	"net/http"
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

	log.Println(r.Context().Value(middleware.UserIDKey))
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
	}

	project, err = pc.projectService.CreateProject(project)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, &utils.ErrorResponse{
			Status:  false,
			Message: "Failed to create project.",
			Errors:  err.Error(),
		})
	}

	utils.WriteJSONResponse(w, http.StatusCreated, &utils.SuccessResponse{
		Status:  true,
		Message: "Project created successfully.",
		Data:    project,
	})
}
