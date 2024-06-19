package routers

import (
	"github.com/drTragger/MykroTask/api/controllers"
	"github.com/gorilla/mux"
)

func SetupRouter(userController *controllers.UserController) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/register", userController.RegisterUser).Methods("POST")
	router.HandleFunc("/api/login", userController.Login).Methods("POST")
	router.HandleFunc("/api/users/{id}", userController.GetUserById).Methods("GET")
	// Define your other routes here

	return router
}
