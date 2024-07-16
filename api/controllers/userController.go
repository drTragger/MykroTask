package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/drTragger/MykroTask/models"
	"github.com/drTragger/MykroTask/services"
	"github.com/drTragger/MykroTask/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (uc *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var userDTO models.CreateUserDTO
	err := json.NewDecoder(r.Body).Decode(&userDTO)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Failed to decode body data",
			Errors:  err.Error(),
		})
		return
	}

	// Validate the request data
	err = utils.ValidateStruct(userDTO)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusUnprocessableEntity, &utils.ErrorResponse{
			Status:  false,
			Message: "Validation failed",
			Errors:  err.Error(),
		})
		return
	}

	if userDTO.Password != userDTO.ConfirmPassword {
		utils.WriteJSONResponse(w, http.StatusUnprocessableEntity, &utils.ErrorResponse{
			Status:  false,
			Message: "Passwords do not match.",
		})
		return
	}

	existingUser, err := uc.userService.GetUserByEmail(userDTO.Email)
	if err == nil && existingUser != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "User with this email already exists.",
		})
		return
	}

	user := models.User{
		Name:      userDTO.Name,
		Email:     userDTO.Email,
		Password:  userDTO.Password,
		CreatedAt: time.Now(),
	}

	err = uc.userService.RegisterUser(&user)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, &utils.ErrorResponse{
			Status:  false,
			Message: "Something went wrong",
			Errors:  err.Error(),
		})
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, &utils.SuccessResponse{
		Status:  true,
		Message: "User successfully registered",
		Data:    user,
	})
	return
}

func (uc *UserController) GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Missing id parameter",
		})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Invalid id parameter",
		})
		return
	}

	user, err := uc.userService.GetUserById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSONResponse(w, http.StatusNotFound, &utils.ErrorResponse{
				Status:  false,
				Message: "User not found",
			})
		} else {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, &utils.ErrorResponse{
				Status:  false,
				Message: "Something went wrong",
				Errors:  err.Error(),
			})
		}
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, &utils.SuccessResponse{
		Status:  true,
		Message: "User found.",
		Data:    user,
	})
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var loginDTO models.LoginDTO
	err := json.NewDecoder(r.Body).Decode(&loginDTO)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Failed to decode body data",
			Errors:  err.Error(),
		})
		return
	}

	// Validate the request data
	err = utils.ValidateStruct(loginDTO)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, &utils.ErrorResponse{
			Status:  false,
			Message: "Validation failed",
			Errors:  err.Error(),
		})
		return
	}

	token, err := uc.userService.LoginUser(loginDTO.Email, loginDTO.Password)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusUnauthorized, &utils.ErrorResponse{
			Status:  false,
			Message: "Invalid email or password",
			Errors:  err.Error(),
		})
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, &utils.SuccessResponse{
		Status:  true,
		Message: "Login successful",
		Data:    token,
	})
	return
}
