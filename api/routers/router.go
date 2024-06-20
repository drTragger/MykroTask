package routers

import (
	"github.com/drTragger/MykroTask/api/controllers"
	"github.com/drTragger/MykroTask/middleware"
	"github.com/gorilla/mux"
)

func SetupRouter(
	userController *controllers.UserController,
	projectController *controllers.ProjectController,
	jwtKey []byte,
) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWTMiddleware(jwtKey))

	// User Management
	router.HandleFunc("/api/register", userController.RegisterUser).Methods("POST")
	router.HandleFunc("/api/login", userController.Login).Methods("POST")
	router.HandleFunc("/api/users/{id}", userController.GetUserById).Methods("GET")

	// Project Management
	api.HandleFunc("/projects", projectController.CreateProject).Methods("POST")
	api.HandleFunc("/projects", projectController.GetProjectsForUser).Methods("GET")
	api.HandleFunc("/projects/{id}", projectController.GetProjectById).Methods("GET")
	api.HandleFunc("/projects/{id}", projectController.UpdateProject).Methods("PUT")

	return router
}
