package controllers

import (
	"github.com/drTragger/MykroTask/middleware"
	"github.com/drTragger/MykroTask/models"
	"github.com/drTragger/MykroTask/services"
	"github.com/drTragger/MykroTask/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

type ProjectMemberController struct {
	projectMemberService services.ProjectMemberService
}

func NewProjectMemberController(projectMemberService services.ProjectMemberService) *ProjectMemberController {
	return &ProjectMemberController{projectMemberService: projectMemberService}
}

func (pmc *ProjectMemberController) CreateMember(w http.ResponseWriter, r *http.Request) {
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

	var member *models.ProjectMember
	errorResponse := utils.UnmarshalRequest(r, &member)
	if errorResponse != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, errorResponse)
		return
	}

	member.ProjectId = projectId
	userId := uuid.MustParse(r.Context().Value(middleware.UserIDKey).(string))

	err = utils.ValidateStruct(member)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Validation failed.",
			Errors:  err.Error(),
		})
		return
	}

	member, forbidden, err := pmc.projectMemberService.CreateMember(member, userId)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, &utils.ErrorResponse{
			Status:  false,
			Message: "Failed to create team member.",
			Errors:  err.Error(),
		})
		return
	}
	if forbidden {
		utils.WriteJSONResponse(w, http.StatusForbidden, &utils.ErrorResponse{
			Status:  false,
			Message: "You are not allowed to create this project member.",
		})
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, &utils.SuccessResponse{
		Status:  true,
		Message: "Team member created successfully.",
		Data:    member,
	})
}

func (pmc *ProjectMemberController) GetMembers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["projectId"]
	if !ok {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Missing projectId parameter.",
		})
		return
	}

	projectId, err := uuid.Parse(idStr)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Invalid projectId parameter.",
		})
		return
	}

	members, err := pmc.projectMemberService.GetMembers(projectId)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, &utils.ErrorResponse{
			Status:  false,
			Message: "Failed to get team members.",
			Errors:  err.Error(),
		})
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, &utils.SuccessResponse{
		Status:  true,
		Message: "Team members retrieved successfully.",
		Data:    members,
	})
}

func (pmc *ProjectMemberController) DeleteMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectIdStr, ok := vars["projectId"]
	if !ok {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Missing projectId parameter.",
		})
		return
	}
	memberIdStr, ok := vars["userId"]
	if !ok {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Missing userId parameter.",
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
	memberId, err := uuid.Parse(memberIdStr)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Invalid userId parameter.",
		})
		return
	}

	userId := uuid.MustParse(r.Context().Value(middleware.UserIDKey).(string))

	forbidden, err := pmc.projectMemberService.DeleteMember(projectId, memberId, userId)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, &utils.ErrorResponse{
			Status:  false,
			Message: "Failed to delete team member.",
			Errors:  err.Error(),
		})
		return
	}
	if forbidden {
		utils.WriteJSONResponse(w, http.StatusForbidden, &utils.ErrorResponse{
			Status:  false,
			Message: "Not allowed to delete team member.",
		})
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, &utils.SuccessResponse{
		Status:  true,
		Message: "Team member deleted successfully.",
	})
}
