package routers

import (
	"github.com/drTragger/MykroTask/api/controllers"
	"github.com/drTragger/MykroTask/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func SetupRouter(
	userController *controllers.UserController,
	projectController *controllers.ProjectController,
	projectMemberController *controllers.ProjectMemberController,
	taskController *controllers.TaskController,
	jwtKey []byte,
) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWTMiddleware(jwtKey))

	// User Management
	router.HandleFunc("/api/register", userController.RegisterUser).Methods(http.MethodPost)
	router.HandleFunc("/api/login", userController.Login).Methods(http.MethodPost)
	router.HandleFunc("/api/users/{id}", userController.GetUserById).Methods(http.MethodGet)

	// Project Management
	api.HandleFunc("/projects", projectController.CreateProject).Methods(http.MethodPost)
	api.HandleFunc("/projects", projectController.GetProjectsForUser).Methods(http.MethodGet)
	api.HandleFunc("/projects/{id}", projectController.GetProjectById).Methods(http.MethodGet)
	api.HandleFunc("/projects/{id}", projectController.UpdateProject).Methods(http.MethodPut)
	api.HandleFunc("/projects/{id}", projectController.DeleteProject).Methods(http.MethodDelete)

	// Project Members Management
	api.HandleFunc("/projects/{projectId}/users", projectMemberController.CreateMember).Methods(http.MethodPost)
	api.HandleFunc("/projects/{projectId}/users", projectMemberController.GetMembers).Methods(http.MethodGet)
	api.HandleFunc("/projects/{projectId}/users/{userId}", projectMemberController.DeleteMember).Methods(http.MethodDelete)

	// Task Management
	api.HandleFunc("/projects/{projectId}/tasks", taskController.CreateTask).Methods(http.MethodPost)
	api.HandleFunc("/projects/{projectId}/tasks", taskController.GetTasksForUser).Methods(http.MethodGet)

	return router
}
