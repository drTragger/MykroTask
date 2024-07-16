package routers

import (
	"github.com/drTragger/MykroTask/api/controllers"
	"github.com/drTragger/MykroTask/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)

func SetupRouter(
	userController *controllers.UserController,
	projectController *controllers.ProjectController,
	projectMemberController *controllers.ProjectMemberController,
	taskController *controllers.TaskController,
	jwtKey []byte,
) http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWTMiddleware(jwtKey))
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Replace with your front-end URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})
	handler := c.Handler(router)

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
	api.HandleFunc("/projects/{projectId}/users/{memberId}/tasks", taskController.GetTasksForUser).Methods(http.MethodGet)
	api.HandleFunc("/projects/{projectId}/tasks/{taskId}", taskController.GetTaskById).Methods(http.MethodGet)
	api.HandleFunc("/projects/{projectId}/tasks/{taskId}", taskController.DeleteTask).Methods(http.MethodDelete)
	api.HandleFunc("/projects/{projectId}/tasks", taskController.GetTasksForProject).Methods(http.MethodGet)
	api.HandleFunc("/projects/{projectId}/tasks/{taskId}", taskController.UpdateTask).Methods(http.MethodPut)

	return handler
}
